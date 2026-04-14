# Cat Breed Helper

Веб-приложение для определения породы кошки по фотографии.

Проект состоит из двух частей:

- **frontend** — интерфейс на HTML, CSS и JavaScript
- **backend** — API на FastAPI

## Demo

- Frontend: `https://catbreed-helper.vercel.app/`
- Backend: `https://catbreed-helper.onrender.com`

## Возможности

- загрузка изображения кошки через веб-интерфейс
- предпросмотр выбранного изображения
- удаление выбранного изображения
- отправка изображения на backend
- отображение результата распознавания
- вывод предполагаемой породы
- вывод уровня уверенности
- отображение источника ответа
- мок-режим для текущей разработки
- готовая структура для подключения реального ML endpoint

## Стек

### Frontend
- HTML5
- CSS3
- JavaScript (Vanilla JS)

### Backend
- FastAPI
- Uvicorn
- python-multipart
- CORS middleware

## Как это работает

1. Пользователь открывает сайт.
2. Выбирает изображение кошки.
3. Интерфейс показывает предпросмотр изображения.
4. После отправки файла frontend делает `POST` запрос на backend.
5. Backend возвращает JSON с результатом.
6. Интерфейс показывает породу, уверенность и источник ответа.

## Структура проекта

```text
catbreed-helper/
├── ai/
├── server/
├── web/
│   ├── backend/
│   │   ├── app.py
│   │   ├── requirements.txt
│   │   └── render.yaml
│   ├── frontend/
│   │   ├── index.html
│   │   ├── styles.css
│   │   └── script.js
│   ├── .gitignore
│   └── README.md
└── README.md
```