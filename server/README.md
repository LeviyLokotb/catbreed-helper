
Запрос: POST на /predict с изоюражением

Ответ: JSON вида:
```json
{
  "breed": "Kitik",
  "confidence": "0.95",
  "filename": "cat_photo_001.jpg"
}
```

Доступна проверка жизнеспособности сервера -- GET запрос на /health