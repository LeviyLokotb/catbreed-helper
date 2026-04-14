import random
from typing import Any

from fastapi import FastAPI, File, UploadFile
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="Cat Breed Helper API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

BREEDS = [
    {"breed_id": "british", "breed": "Британская короткошёрстная"},
    {"breed_id": "maine_coon", "breed": "Мейн-кун"},
    {"breed_id": "siamese", "breed": "Сиамская"},
    {"breed_id": "sphynx", "breed": "Сфинкс"},
    {"breed_id": "scottish_fold", "breed": "Шотландская вислоухая"},
]


@app.get("/health")
def healthcheck() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/predict")
async def predict(file: UploadFile = File(...)) -> dict[str, Any]:
    if not file.content_type or not file.content_type.startswith("image/"):
        return {
            "success": False,
            "error": "Нужно загрузить изображение.",
        }

    contents = await file.read()
    if not contents:
        return {
            "success": False,
            "error": "Файл пустой.",
        }

    prediction = random.choice(BREEDS)
    confidence = round(random.uniform(0.82, 0.98), 4)

    return {
        "success": True,
        "breed_id": prediction["breed_id"],
        "breed": prediction["breed"],
        "confidence": confidence,
        "filename": file.filename,
        "mock": True,
    }