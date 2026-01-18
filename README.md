# subscriptions-service

REST-сервис для управления подписками пользователей и подсчёта суммарной стоимости подписок за период.

## Стек
- Go (chi, pgx)
- PostgreSQL
- migrate/migrate (SQL migrations)
- OpenAPI 3.0 (Swagger UI)

## Возможности
- CRUDL над подписками:
  - Create `POST /api/v1/subscriptions`
  - Read `GET /api/v1/subscriptions/{id}`
  - Update `PUT /api/v1/subscriptions/{id}`
  - Delete `DELETE /api/v1/subscriptions/{id}`
  - List `GET /api/v1/subscriptions`
- Подсчёт total за период:
  - `GET /api/v1/subscriptions/total`
- Формат дат в API: `MM-YYYY` (например `07-2025`)
- Валидация `user_id` как UUID

## Требования
- Docker + Docker Compose
- (опционально) Go для локального запуска

---

## Запуск через Docker Compose

### 1) Создать `.env`
```bash
cp .env.example .env
````

### 2) Собрать и поднять сервисы

```bash
docker compose up -d --build
```

### 3) Применить миграции

```bash
docker compose run --rm migrate
```

### 4) Проверка жив ли сервис

```bash
curl -i http://localhost:8080/health
curl -i http://localhost:8080/db/health
```

---

## Swagger / OpenAPI

* Swagger UI: `http://localhost:8081`
* Спецификация: `api/openapi.yaml`

---

## Переменные окружения

* `APP_PORT` — порт приложения (по умолчанию 8080)
* `DB_HOST` — хост Postgres (в docker compose: `postgres`)
* `DB_PORT` — порт Postgres (в docker compose: `5432`)
* `DB_USER` — пользователь
* `DB_PASSWORD` — пароль
* `DB_NAME` — база данных
* `DB_SSLMODE` — режим ssl

---

## Примеры запросов

### Create

#### Без end_date

```bash
curl -i -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 99,
    "user_id": "11111111-1111-1111-1111-111111111111",
    "start_date": "07-2025",
    "end_date": null
  }'
```

#### С end_date

```bash
curl -i -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 199,
    "user_id": "11111111-1111-1111-1111-111111111111",
    "start_date": "07-2025",
    "end_date": "10-2025"
  }'
```

### Read (Get by id)

```bash
curl -i http://localhost:8080/api/v1/subscriptions/<ID>
```

### List

```bash
curl -i http://localhost:8080/api/v1/subscriptions
curl -i "http://localhost:8080/api/v1/subscriptions?user_id=11111111-1111-1111-1111-111111111111"
curl -i "http://localhost:8080/api/v1/subscriptions?service_name=Yandex%20Plus"
curl -i "http://localhost:8080/api/v1/subscriptions?user_id=11111111-1111-1111-1111-111111111111&service_name=Yandex%20Plus"
```

### Update

```bash
curl -i -X PUT http://localhost:8080/api/v1/subscriptions/<ID> \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 500,
    "user_id": "11111111-1111-1111-1111-111111111111",
    "start_date": "07-2025",
    "end_date": "10-2025"
  }'
```

### Delete

```bash
curl -i -X DELETE http://localhost:8080/api/v1/subscriptions/<ID>
```

### Total

```bash
curl -i "http://localhost:8080/api/v1/subscriptions/total?user_id=11111111-1111-1111-1111-111111111111&start_date=07-2025&end_date=10-2025"
curl -i "http://localhost:8080/api/v1/subscriptions/total?user_id=11111111-1111-1111-1111-111111111111&service_name=Yandex%20Plus&start_date=07-2025&end_date=10-2025"
```

---

## Остановка

```bash
docker compose down
```

Удалить данные Postgres:

```bash
docker compose down -v
```