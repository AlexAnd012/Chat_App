FROM golang:1.24-alpine AS build
WORKDIR /app

# Ускоряем кеширование модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальное и собираем
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd

FROM alpine:3.20
RUN adduser -D -H appuser
WORKDIR /app

# Только бинарник
COPY --from=build /app/server /app/server

# Переменные среды по умолчанию 
ENV ADDR=":8080"
EXPOSE 8080

USER appuser
ENTRYPOINT ["/app/server"]
