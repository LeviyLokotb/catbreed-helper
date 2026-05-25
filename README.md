# Cat Breed Helper

Приложение для определения породы кошки по фото с реальным ML-инференсом (без мока).

## Структура проекта

- `ai/` — ноутбуки и артефакты обучения/экспорта модели.
- `server/` — Go ML-сервис (ONNX Runtime), эндпоинты `/health` и `/predict`.
- `web/backend/` — FastAPI backend-прокси для фронта.
- `web/frontend/` — HTML/CSS/JS интерфейс.

## Как работает цепочка

1. Пользователь загружает или вставляет изображение на фронте.
2. Frontend отправляет файл в `web/backend` (`POST /predict`).
3. Backend отправляет файл в ML `server` (`POST /predict`).
4. ML возвращает породу и confidence, frontend показывает результат.

## Локальный запуск

Требования:
- Docker Desktop
- Python 3.11+ (для `web/backend` и локальной раздачи фронта)

### 1) Запуск ML server (Go + ONNX) в Docker

```powershell
cd d:\catbreed-helper\server
docker compose down
docker compose up --build
```

Проверка:
- `http://127.0.0.1:8081/health` -> `OK`

### 2) Запуск backend (FastAPI)

```powershell
cd d:\catbreed-helper\web\backend
python -m venv .venv
.\.venv\Scripts\activate
pip install -r requirements.txt
$env:LOCAL_DEV="true"
uvicorn app:app --host 0.0.0.0 --port 8000 --reload
```

Проверка:
- `http://127.0.0.1:8000/health` -> `{"status":"ok"}`
- `http://127.0.0.1:8000/` -> должен показывать `ml_predict_url` на `http://127.0.0.1:8081/predict`

### 3) Запуск frontend

```powershell
cd d:\catbreed-helper\web\frontend
python -m http.server 5500
```

Открыть в браузере:
- `http://127.0.0.1:5500`

## Поддерживаемые форматы

- JPG
- PNG

Также поддерживается вставка картинки из буфера обмена (`Ctrl+V`).

## Что важно по UX

- При низкой уверенности модель возвращает ошибку "не удалось уверенно распознать кошку".
- Ошибки на фронте подсвечиваются красной плашкой.
