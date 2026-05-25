FROM golang:1.25-bookworm AS builder

ENV GOPROXY=direct

# Устанавливаем зависимости для сборки
RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Скачиваем и устанавливаем ONNX Runtime
RUN wget https://github.com/microsoft/onnxruntime/releases/download/v1.25.0/onnxruntime-linux-x64-1.25.0.tgz \
    && tar -xzf onnxruntime-linux-x64-1.25.0.tgz \
    && cp onnxruntime-linux-x64-1.25.0/lib/libonnxruntime* /usr/lib/ \
    && cp -r onnxruntime-linux-x64-1.25.0/include/* /usr/include/ \
    && ln -s /usr/lib/libonnxruntime.so.1.25.0 /usr/lib/onnxruntime.so \
    && rm -rf onnxruntime-linux-x64-1.25.0*

WORKDIR /app

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем сервер
RUN go build -o server ./cmd/main/

# Финальный образ
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Копируем ONNX Runtime
COPY --from=builder /usr/lib/libonnxruntime* /usr/lib/
COPY --from=builder /usr/lib/onnxruntime.so /usr/lib/
RUN ldconfig

# Копируем сервер
COPY --from=builder /app/server /app/server

# Копируем модель и labels
COPY model.onnx /app/model.onnx
COPY labels.txt /app/labels.txt

ENV MODEL_PATH=/app/model.onnx
ENV LABELS_PATH=/app/labels.txt

EXPOSE 8080
CMD ["/app/server"]