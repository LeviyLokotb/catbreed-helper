from __future__ import annotations

import os
from typing import Any

import httpx
from fastapi import FastAPI, File, HTTPException, UploadFile
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="Cat Breed Helper API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Для локальной разработки можно принудительно использовать локальный ML:
# LOCAL_DEV=true -> http://127.0.0.1:8081/predict
LOCAL_DEV = os.getenv("LOCAL_DEV", "true").strip().lower() in {"1", "true", "yes", "on"}
DEFAULT_LOCAL_ML_URL = "http://127.0.0.1:8081/predict"
ML_PREDICT_URL = (
    DEFAULT_LOCAL_ML_URL
    if LOCAL_DEV
    else os.getenv("ML_PREDICT_URL", DEFAULT_LOCAL_ML_URL).strip()
)
REQUEST_TIMEOUT_SECONDS = float(os.getenv("ML_REQUEST_TIMEOUT_SECONDS", "30"))
MIN_CAT_CONFIDENCE = float(os.getenv("MIN_CAT_CONFIDENCE", "0.40"))


@app.get("/")
def root() -> dict[str, Any]:
    return {
        "ok": True,
        "message": "Cat Breed Helper API is running",
        "ml_predict_url_configured": bool(ML_PREDICT_URL),
        "local_dev": LOCAL_DEV,
        "ml_predict_url": ML_PREDICT_URL,
    }


@app.get("/health")
def healthcheck() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/predict")
async def predict(file: UploadFile = File(...)) -> dict[str, str]:
    allowed_types = {"image/jpeg", "image/png"}
    if not file.content_type or file.content_type.lower() not in allowed_types:
        raise HTTPException(status_code=400, detail="Поддерживаются только изображения JPG и PNG.")

    file_bytes = await file.read()
    if not file_bytes:
        raise HTTPException(status_code=400, detail="Файл пустой.")

    if not ML_PREDICT_URL:
        raise HTTPException(
            status_code=500,
            detail="Не задан ML_PREDICT_URL. Укажи URL реальной нейронки в переменной окружения.",
        )

    filename = file.filename or "image.jpg"
    content_type = file.content_type or "application/octet-stream"

    try:
        async with httpx.AsyncClient(timeout=REQUEST_TIMEOUT_SECONDS, trust_env=False) as client:
            response = await client.post(
                ML_PREDICT_URL,
                files={"file": (filename, file_bytes, content_type)},
            )
        response.raise_for_status()
    except httpx.TimeoutException as exc:
        raise HTTPException(status_code=504, detail="ML сервис не ответил вовремя.") from exc
    except httpx.HTTPStatusError as exc:
        upstream_preview = exc.response.text[:200] if exc.response is not None else ""
        raise HTTPException(
            status_code=502,
            detail=f"ML сервис вернул ошибку: {exc.response.status_code}. {upstream_preview}",
        ) from exc
    except httpx.RequestError as exc:
        raise HTTPException(status_code=502, detail="Не удалось связаться с ML сервисом.") from exc

    try:
        payload = response.json()
    except ValueError as exc:
        raise HTTPException(status_code=502, detail="ML сервис вернул невалидный JSON.") from exc

    breed = str(payload.get("breed", "")).strip()
    confidence = str(payload.get("confidence", "")).strip()
    returned_filename = str(payload.get("filename", filename)).strip()

    if not breed or not confidence:
        raise HTTPException(
            status_code=502,
            detail="Ответ ML сервиса не содержит обязательные поля breed/confidence.",
        )

    try:
        confidence_value = float(confidence)
    except ValueError as exc:
        raise HTTPException(status_code=502, detail="Невалидное поле confidence в ответе ML сервиса.") from exc

    if confidence_value < MIN_CAT_CONFIDENCE:
        raise HTTPException(
            status_code=422,
            detail=(
                "Не удалось уверенно распознать кошку на изображении. "
                "Попробуйте фото, где кошка крупнее и лучше освещена."
            ),
        )

    return {
        "breed": breed,
        "confidence": str(confidence_value),
        "filename": returned_filename,
    }
