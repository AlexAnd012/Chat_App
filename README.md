# Chat App (Go) — Backend
Мини-чат на Go с хранением истории в Redis, сборкой и запуском через Docker/Compose, и интерактивной документацией API на Swagger.

## **Стек и как используется**
- **Go (Gin)** — HTTP-сервер, маршрутизация и обработчики.

- **Redis** — быстрая лента сообщений:

- **PING** — health-check соединения.

- **Docker** / Docker Compose — воспроизводимая сборка и запуск (api + redis) одной командой.

- **Swagger (swag + gin-swagger)** — живой UI на /swagger/index.html с кнопкой **Try it out**. Спецификация генерируется из комментариев в коде.

## **Реализовано**

- **POST /api/messages** — создать сообщение (валидация username/text, запись в Redis).

- **GET /api/messages?limit=N** — получить последние N сообщений в хронологическом порядке.

- **GET /healthz** — проверка живости сервиса и доступности Redis.

- **Swagger UI по пути /swagger/index.html.**

- **Dockerfile (multi-stage) и docker-compose.yml** для локального запуска.

## **старт**
**docker compose up --build**
Доступ по умолчанию: http://localhost:8080.

**Проверка:**
curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/api/messages \
  -H 'Content-Type: application/json' \
  -d '{"username":"Alex","text":"Hello from Docker!"}'
curl 'http://localhost:8080/api/messages?limit=10'

**Открыть Swagger UI:**
http://localhost:8080/swagger/index.html

## **Переменные окружения**
(из .env в docker-compose.yml)

ADDR=:8080
REDIS_ADDR=redis:6379     # в докере имя сервиса redis
HISTORY_LIMIT=500         # сколько сообщений храним/отдаём максимум
GIN_MODE=release

Для локального запуска без Docker: REDIS_ADDR=127.0.0.1:6379.

## **API** 

GET /healthz → 200 OK, если сервис и Redis доступны.

POST /api/messages
{ "username": "Alex", "text": "Hello!" }

Ответ 201 Created:
{
  "id": "uuid",
  "username": "Alice",
  "text": "Hello!",
  "ts": "2025-01-01T12:00:00Z"
}

GET /api/messages?limit=50
Возвращает массив последних сообщений (старые → новые).

## **Структура проекта**
.
├── cmd/
│   └── server/
│       └── main.go           # запуск HTTP-сервера, DI зависимостей
├── internal/
│   ├── httpserver/
│   │   ├── handlers.go       # POST/GET messages, валидация, ответы
│   │   └── router.go         # маршруты (/healthz, /api/*, /swagger)
│   └── storage/
│       └── redisrepo.go      # Redis-клиент: Ping, AppendMessage, TrimHistory, RecentMessages
├── docs/                     # сгенерированные Swagger-файлы (swag init)
├── go.mod / go.sum
├── Dockerfile                # multi-stage сборка
├── docker-compose.yml        # api + redis
└── .env                      # конфиг окружения
