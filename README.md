
<div align="center">
  <h1>🧮 Распределённый вычислитель 🧮
<br> арифметических выражений </h1>
<h4>персистентность и многопользовательский режим</h4>
</div>


## <calc final super pooper 5>

### Содержание 📌

- [ℹ️ О проекте](#ℹ️-о-проекте)
- [⬇️ Установка](%EF%B8%8F-установка)
- [▶️ Запуск](#▶️-запуск)
- [⁉️ Примеры использования](#⁉️-примеры-использования) 
- [🎊 Тесты](#запуск-тестов) 
- [📖 Документация "Как это всё работает?"](#📖-документация)

## ℹ️ О проекте 
Это веб-сервис, который позволяет пользователям быстро, а иногда и не очень вычислять арифметические выражения. Пользователь отправляет арифметическое выражение через HTTP-запрос, и в ответ он получает результат вычисления. 

За основу взято старое решение веб-сервиса, но его функционал полностью переписан. Новое решение вычисляет части арифметического выражения параллельно, что значительно ускоряет процесс.

## ⬇️ Установка
Запустить проект можно несколькими способами. [Далее](#запуск) мы разберем несколько из них.

Но для начала, нужно установить сервис на ваше устройство. Для этого клонируйте репозиторий.
```bash
git clone https://github.com/sklerakuku/5final.git
```
Перейдите в директорию проекта
```bash
cd 5final
```
## ▶️ Запуск
Запустите сервисы оркестратора и агента поочередно в разных терминалах.
```bash
go run ./cmd/orchestrator/main.go
```
```bash
go run ./cmd/agent/main.go
```
Или же одной командой
```bash
go run ./cmd/orchestrator/main.go && go run ./cmd/agent/main.go
```
Всё готово!! Теперь вы сервис доступен по адресу http://localhost:8080/

> [!TIP]
> P.S.: Для остановки сервера нажмите сочетание клавиш: `<CTRL-C>`

## ⁉️ Примеры использования
Вы можете воспользоваться для запросов как curl-ом, так и Postman-ом. Или использовать пользовательский интерфейс веб-сервиса.

![Как пользоваться веб-интерфейсом](https://github.com/sklerakuku/5final/blob/main/images/howtouse.png)

> _login: **admin** password: **admin123**_

 **Переменные  среды сервиса**

*DATABASE_PATH - путь к файлу базы SQLite (default: "./calculator.db")
SERVER_PORT - порт HTTP сервера (default: "8080")
JWT_SECRET - Secret for JWT tokens (default: "secret")
GRPC_ADDRESS - адресс gRPC сервера (default: "localhost:50051")
COMPUTING_POWER - количество горутин
TIME_ADDITION_MS - время выполнения операции сложения в миллисекундах  
TIME_SUBTRACTION_MS - время выполнения операции вычитания в миллисекундах  
TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в миллисекундах  
TIME_DIVISIONS_MS - время выполнения операции деления в миллисекундах*

# 🧮 API Документация калькулятора выражений

## 🔐 Аутентификация
Все эндпоинты, кроме `/api/v1/register` и `/api/v1/login`, требуют JWT токена в заголовке:
```
Authorization: Bearer <ваш_токен>
```

---

## 📝 Эндпоинты API

### 🔹 Регистрация нового пользователя
`POST /api/v1/register`

**Тело запроса:**
```json
{
    "login": "username",
    "password": "password123"
}
```

**Требования:**
- Логин: минимум 3 символа
- Пароль: минимум 6 символов

**Пример:**
```bash
curl -X POST http://localhost:8080/api/v1/register \
-H "Content-Type: application/json" \
-d '{"login":"testuser", "password":"secret123"}'
```

**Успешный ответ (200 OK):**
```json
{"message": "Registration successful"}
```

**Ошибки:**
- `400 Bad Request` - Неверный формат запроса
- `409 Conflict` - Пользователь уже существует
- `500 Internal Server Error` - Ошибка сервера

---

### 🔹 Авторизация
`POST /api/v1/login`

**Тело запроса:**
```json
{
    "login": "username",
    "password": "password123"
}
```

**Пример:**
```bash
curl -X POST http://localhost:8080/api/v1/login \
-H "Content-Type: application/json" \
-d '{"login":"testuser", "password":"secret123"}'
```

**Успешный ответ (200 OK):**
```json
{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
```

**Ошибки:**
- `400 Bad Request` - Неверный формат запроса
- `401 Unauthorized` - Неверные учетные данные

---

### 🔹 Вычисление выражения
`POST /api/v1/calculate`

**Требует аутентификации**

**Тело запроса:**
```json
{
    "expression": "2+2*2"
}
```

**Поддерживаемые операции:**
- `+` - сложение
- `-` - вычитание
- `*` - умножение
- `/` - деление
- `()` - скобки

**Пример запроса:**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
-d '{"expression": "2+2*2"}'
```

**Успешный ответ (200 OK):**
```json
{
    "id": 42,
    "result": 6
}
```

**Примеры ошибок:**

1. Неверное выражение (422 Unprocessable Entity):
```json
{"message": "invalid number format: abc"}
```

2. Деление на ноль (400 Bad Request):
```json
{"message": "division by zero"}
```

3. Ошибка сервера (500 Internal Server Error):
```json
{"message": "Internal server error"}
```

---

### 🔹 Получение выражения по ID
`GET /api/v1/expressions/{id}`

**Требует аутентификации**

**Пример запроса:**
```bash
curl -X GET http://localhost:8080/api/v1/expressions/42 \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Успешный ответ (200 OK):**
```json
{
    "id": 42,
    "user_id": 1,
    "expression": "2+2*2",
    "status": "completed",
    "result": 6
}
```

**Ошибки:**
- `404 Not Found` - Выражение не найдено
- `403 Forbidden` - Нет доступа к выражению

---

### 🔹 Получение всех выражений пользователя
`GET /api/v1/expressions`

**Требует аутентификации**

**Пример запроса:**
```bash
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Успешный ответ (200 OK):**
```json
{
    "expressions": [
        {
            "id": 42,
            "expression": "2+2*2",
            "status": "completed",
            "result": 6
        },
        {
            "id": 43,
            "expression": "10/0",
            "status": "failed",
            "result": null
        }
    ]
}
```

---

## 🚀 Примеры использования

### 1. Вычисление простого выражения
**Запрос:**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGci..." \
-d '{"expression": "(3+5)*2"}'
```

**Ответ:**
```json
{"id": 44, "result": 16}
```

### 2. Ошибка в выражении
**Запрос:**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGci..." \
-d '{"expression": "10/0"}'
```

**Ответ:**
```json
{"message": "division by zero"}
```

### 3. Получение истории вычислений
**Запрос:**
```bash
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Authorization: Bearer eyJhbGci..."
```

**Ответ:**
```json
{
    "expressions": [
        {"id": 44, "expression": "(3+5)*2", "status": "completed", "result": 16},
        {"id": 45, "expression": "2^3", "status": "failed", "result": null}
    ]
}
```

## Запуск тестов

Для запуска тестов используйте команды:

```bash
# Модульные тесты
go test -v ./internal/calculation/...

# Интеграционные тесты
go test -v ./integration/...
```


## 📖  Документация
> "*КАК ВСЁ ЭТО РАБОТАЕТ?*"

![Упрощённая схема](https://github.com/sklerakuku/calc3.0/blob/bfc3483c6f73fceaa8c07843b1c9560b2f48e740/123%20(1).png)

<details>
  <summary>Подробная схема</summary>

```mermaid
sequenceDiagram
    participant Client as HTTP Client
    participant Auth as Auth Service
    participant Server as HTTP Server
    participant DB as Database
    participant Parser as Expression Parser
    participant gRPC as gRPC Client
    participant Worker as gRPC Worker

    Client->>Auth: POST /register (login, password)
    Auth->>DB: Save user
    DB-->>Auth: User created
    Auth-->>Client: Success

    Client->>Auth: POST /login (login, password)
    Auth->>DB: Verify credentials
    DB-->>Auth: User data
    Auth-->>Client: JWT token

    Client->>Server: POST /calculate (expression) + JWT
    Server->>Auth: Verify JWT
    Auth-->>Server: Valid claims
    Server->>DB: Save expression
    DB-->>Server: Expression ID
    Server->>Parser: Parse expression
    Parser-->>Server: Parsed operations
    Server->>gRPC: Send operations
    gRPC->>Worker: Distribute tasks
    Worker-->>gRPC: Results
    gRPC-->>Server: Aggregated result
    Server->>DB: Update expression status
    DB-->>Server: Success
    Server-->>Client: Expression ID
```
</details>

## Общее описание системы

Это система для вычисления математических выражений с аутентификацией пользователей и распределенной обработкой через gRPC. Основные компоненты:

1. **Сервер аутентификации** - регистрация и авторизация пользователей (JWT)
2. **Парсер выражений** - разбор математических выражений
3. **gRPC сервис** - распределенные вычисления
4. **База данных** - хранение пользователей и выражений
5. **HTTP сервер** - API для взаимодействия с клиентами


----
> [!TIP]
> Привет проверяющий!
> 😊



