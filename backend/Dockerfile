# backend/Dockerfile
FROM golang:latest

WORKDIR /app

# Сначала копируем файлы зависимостей и скачиваем их
# Это ускоряет сборку, если код поменялся, а библиотеки нет
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь остальной код
COPY . .

# Собираем бинарник из твоей точки входа
RUN go build -o main ./cmd/calculator/main.go

EXPOSE 8080

# Запускаем приложение
CMD ["./main"]