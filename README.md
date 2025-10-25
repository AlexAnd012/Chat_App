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

curl http://localhost:8080/healthz <br> curl -X POST http://localhost:8080/api/messages \ <br> 
-H 'Content-Type: application/json' \ <br>-d '{"username":"Alex","text":"Hello from Docker!"}' <br> curl 'http://localhost:8080/api/messages?limit=10'

**Открыть Swagger UI:**
http://localhost:8080/swagger/index.html

## **Переменные окружения**
**(из .env в docker-compose.yml)**

- ADDR=:8080<br>
- REDIS_ADDR=redis:6379     # в докере имя сервиса redis<br>
- HISTORY_LIMIT=500         # сколько сообщений храним/отдаём максимум<br>
- GIN_MODE=release<br>

## **API** 

- **GET /healthz** → 200 OK, если сервис и Redis доступны.

- **POST /api/messages**
{ "username": "Alex", "text": "Hello!" }

Ответ 201 Created:
{
  "id": "uuid",
  "username": "Alice",
  "text": "Hello!",
  "ts": "2025-01-01T12:00:00Z"
}

- **GET /api/messages?limit=50**
Возвращает массив последних сообщений (старые → новые).

## **Структура проекта**
.<br>
├── cmd/<br>
│   └── main.go   # запуск HTTP-сервера, DI зависимостей<br>   
├── internal/<br>
│   ├── data/<br>
│   │   └── types.go <br>
│   ├── httpserver/<br>
│   │   ├── handlers.go       # POST/GET messages, валидация, ответы<br>
│   │   └── router.go         # маршруты (/healthz, /api/*, /swagger)<br>
│   └── storage/<br>
│       └── Redisrepo.go      # Redis-клиент: Ping, AppendMessage, TrimHistory, RecentMessages<br>
├── docs/                     # сгенерированные Swagger-файлы (swag init)<br>
├── go.mod / go.sum<br>
├── Dockerfile                # multi-stage сборка<br>
├── docker-compose.yml        # api + redis<br>
└── .env                      # конфиг окружения<br>
