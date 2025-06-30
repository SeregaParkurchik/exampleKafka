# Dockerfile для app-producer

# Используем образ Go для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера.
# Это будет корневая директория, куда мы скопируем весь проект.
WORKDIR /app/activity_tracker

# Копируем go.mod и go.sum из корня проекта в /app/activity_tracker
COPY go.mod .
COPY go.sum .

# Загружаем все Go-модули
RUN go mod download

# Копируем только необходимые директории с исходным кодом
# Эти пути `app-producer/`, `config/` указывают на директории внутри контекста сборки (который является корнем вашего проекта)
COPY app-producer/ ./app-producer/
COPY config/ ./config/ 
COPY pkg/ ./pkg/

# Собираем приложение
# -o /app/producer: указывает имя выходного исполняемого файла (вне текущего WORKDIR для удобства)
# ./app-producer/cmd/producer: теперь путь корректен относительно текущего WORKDIR /app/activity_tracker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/producer ./app-producer/cmd/producer

# Используем легковесный образ для финального образа
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный исполняемый файл из образа-строителя
COPY --from=builder /app/producer .

# Открываем порт, на котором слушает приложение
EXPOSE 8080

# Команда для запуска приложения при старте контейнера
CMD ["./producer"]