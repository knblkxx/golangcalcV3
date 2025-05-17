# Калькулятор 
Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения и получать результаты их вычисления.
Асинхронный HTTP-сервис для вычисления арифметических выражений с поддержкой приоритетов операций и скобок. Результаты запросов сохраняются в SQLite.

---

## Возможности

* **Асинхронная обработка:** отправляете POST-запрос с выражением — получаете `id`, по которому позже можно запросить результат.
* **Приоритет операций и скобки:** поддерживаются `+`, `-`, `*`, `/` с корректным порядком вычисления и вложенными скобками.
* **Внутренний оркестратор:** операции делятся на задачи, отправляемые на внутренний endpoint `/internal/task`.
* **Хранение результатов:** SQLite-база хранит поля `id`, `status` и `result`.

---

## Требования

* Go 1.18 или выше
* Git

---

## Установка и запуск

1. **Клонировать репозиторий**:

   ```bash
   git clone https://github.com/knblkxx/golangcalcV3.git
   cd golangcalcV3
   ```

2. **Установить зависимости**:

   ```bash
   go mod tidy
   ```

3. **Настроить переменные среды** (опционально) в файле `.env`:

   ```ini
   COMPUTING_POWER=1000         # макс. параллельных задач
   TIME_ADDITION_MS=0           # задержка сложения в мс
   TIME_SUBTRACTION_MS=0        # задержка вычитания в мс
   TIME_MULTIPLICATIONS_MS=0    # задержка умножения в мс
   TIME_DIVISIONS_MS=0          # задержка деления в мс
   ```

4. **Инициализация БД** (скрипт выполнится автоматически при старте):

   * Таблица `auth` для пользователей
   * Таблица `expressions` для запросов

5. **Запустить сервер**:

   ```bash
   go run main/main.go
   ```

   По умолчанию слушает порт `8080`.

---

## API

### Регистрация пользователя

* **URL:** `POST /api/v1/register`
* **Body:**

  ```json
  {
    "login": "user1",
    "password": "password123"
  }
  ```
* **Response:**

  ```json
  { "status": "registered" }
  ```

### Вход (получение JWT)

* **URL:** `POST /api/v1/login`
* **Body:**

  ```json
  {
    "login": "user1",
    "password": "password123"
  }
  ```
* **Response:**

  ```json
  {
    "token": "<JWT токен>"
  }
  ```

### Отправка выражения на вычисление

* **URL:** `POST /api/v1/calculate`
* **Заголовки:**

  * `Content-Type: application/json`
  * `Authorization: OAuth <JWT токен>`
* **Body:**

  ```json
  { "expression": "2+2*2" }
  ```
* **Response:**

  ```json
  { "id": 42 }
  ```

### Получение статуса и результата

* **URL:** `GET /api/v1/expressions/{id}`
* **Заголовки:**

  * `Authorization: OAuth <JWT токен>`
* **Responses:**

  * **В процессе**

    ```json
    { "id": 42, "status": "proccessing", "result": 0 }
    ```
  * **Успешно**

    ```json
    { "id": 42, "status": "ok", "result": 6 }
    ```
  * **Ошибка**

    ```json
    { "id": 42, "status": "error", "result": 0 }
    ```

### Получение списка запросов

* **URL:** `GET /api/v1/expressions`
* **Заголовки:**

  * `Authorization: OAuth <JWT токен>`
* **Response:**

  ```json
  {
    "expressions": [
      { "id": 42, "status": "ok", "result": 6 },
      { "id": 43, "status": "processing", "result": 0 }
    ]
  }
  ```

### Примеры cURL
Не забудьте открыть git bash !

1. **Регистрация**:

   ```bash
   curl -X POST http://localhost:8080/api/v1/register \
     -H 'Content-Type: application/json' \
     -d '{"login":"user1","password":"pass"}'
   ```

2. **Логин и получение токена (без `jq`)**:

   ```bash
   # Запрос логина и сохранение полного ответа в переменную
   RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
     -H 'Content-Type: application/json' \
     -d '{"login":"user1","password":"pass"}')
   # Извлечение поля token регулярным выражением
   TOKEN=$(echo "$RESPONSE" | grep -oP '(?<="token":")[^" ]*')
   echo "JWT: $TOKEN"
   ```

3. **Отправка выражения**:

   ```bash
   curl -X POST http://localhost:8080/api/v1/calculate \
     -H 'Content-Type: application/json' \
     -H "Authorization: OAuth $TOKEN" \
     -d '{"expression":"(1+2)*3"}'
   ```

4. **Получение результата**:

   ```bash
   curl http://localhost:8080/api/v1/expressions/42 \
     -H "Authorization: OAuth $TOKEN"
   ```

---

## Структура проекта

```
├── auth          # Регистрация, логин и JWT middleware
├── calculator    # Логика парсинга и вычисления выражений
├── handlers      # HTTP‑обработчики
├── main          # Инициализация и маршруты
├── database      # SQLite-файл и схема
├── structs       # Общие типы запросов/ответов
└── README.md     # Документация
```

---

## Тестирование

Запуск модульных и интеграционных тестов:

```bash
go test ./calculator ./handlers ./auth
```

---
