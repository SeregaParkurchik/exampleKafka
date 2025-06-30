# Dockerfile для app-consumer

# Используем образ Go для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера.
WORKDIR /app/activity_tracker

# Копируем go.mod и go.sum из корня проекта в /app/activity_tracker
COPY go.mod .
COPY go.sum .

# Загружаем все Go-модули
RUN go mod download

# Копируем только необходимые директории с исходным кодом
# Эти пути `app-consumer/`, `config/` указывают на директории внутри контекста сборки
COPY app-consumer/ ./app-consumer/
COPY config/ ./config/ 
COPY pkg/ ./pkg/

# Собираем приложение
# -o /app/consumer: указывает имя выходного исполняемого файла (вне текущего WORKDIR для удобства)
# ./app-consumer/cmd/consumer: теперь путь корректен относительно текущего WORKDIR /app/activity_tracker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/consumer ./app-consumer/cmd/consumer

# Используем легковесный образ для финального образа
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный исполняемый файл из образа-строителя
COPY --from=builder /app/consumer .

# Команда для запуска приложения при старте контейнера
CMD ["./consumer"]