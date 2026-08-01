# To Do App (Golang)

Простой REST API для управления задачами, написанный на Go. Приложение построено по принципу "feature-based" архитектуры: каждая доменная область (`users`, `tasks`, `statistics`) содержит свои слои `service` / `repository` / `transport`, а общий код (логгер, HTTP-сервер, работа с Postgres) вынесен в `internal/core`.

## Стек технологий

- **Go** 1.26
- **PostgreSQL 18** — хранилище данных, доступ через [pgx/v5](https://github.com/jackc/pgx)
- **golang-migrate** — миграции БД
- **zap** — структурированное логирование
- **go-playground/validator** — валидация входных данных
- **net/http** (`http.ServeMux`) — HTTP-роутинг, без сторонних фреймворков
- **Docker Compose** — окружение для Postgres и миграций

## Структура проекта

```
cmd/todoapp/          — точка входа (main.go)
internal/core/         — переиспользуемая инфраструктура
  ├── domain/          — доменные сущности (User, Task, Statistics) и их валидация
  ├── errors/          — общие типы ошибок
  ├── logger/          — абстракция логгера + реализация на zap
  ├── repository/      — пул соединений с Postgres (pgx)
  └── transport/http/  — HTTP-сервер, роутер, middleware, request/response утилиты
internal/features/
  ├── users/           — CRUD пользователей
  ├── tasks/           — CRUD задач
  └── statistics/      — статистика
migrations/            — SQL-миграции (golang-migrate)
docker-compose.yaml    — Postgres + туннель для локального порта + сервис миграций
Makefile               — команды для запуска окружения, миграций и приложения
```

## Настройка окружения

1. Скопируйте файл с переменными окружения и заполните его:

   ```bash
   cp .env.example .env
   ```

   Основные переменные:

   | Переменная            | Описание                                   | По умолчанию |
   |------------------------|---------------------------------------------|--------------|
   | `POSTGRES_USER`        | пользователь Postgres                       | —            |
   | `POSTGRES_PASSWORD`    | пароль Postgres                             | —            |
   | `POSTGRES_DB`          | имя базы данных                             | —            |
   | `POSTGRES_TIMEOUT`     | таймаут подключения к БД                    | `10s`        |
   | `HTTP_ADDR`            | адрес, на котором слушает HTTP-сервер       | `:5050`      |
   | `HTTP_SHUTDOWN_TIMEOUT`| таймаут graceful shutdown сервера           | `30s`        |
   | `LOGGER_LEVEL`         | уровень логирования    | `DEBUG`      |

## Запуск

Все основные операции вынесены в `Makefile`.

1. Поднять Postgres и пробросить порт на `127.0.0.1:5432`:

   ```bash
   make postgres
   ```

   (эквивалент `make env-up` + `make env-port-forward`)

2. Применить миграции:

   ```bash
   make migrate-up
   ```

3. Запустить приложение:

   ```bash
   make todoapp-run
   ```

   По умолчанию сервер поднимется на `http://localhost:5050`, логи пишутся в `out/logs/`.

### Остановка и очистка окружения

```bash
make env-down          # остановить контейнер Postgres
make env-port-close    # закрыть проброс порта
make env-cleanup       # удалить volume с данными Postgres (с подтверждением)
make logs-cleanup      # удалить локальные логи (с подтверждением)
```

### Работа с миграциями

```bash
make migrate-create seq=<название>   # создать новую миграцию
make migrate-up                       # применить миграции
make migrate-down                     # откатить миграции
make migrate-action action=<action>   # произвольное действие migrate (up/down/version/force)
```

## API

Все эндпоинты доступны с префиксом `/api/v1`.

### Users

| Метод  | Путь              | Описание                |
|--------|-------------------|--------------------------|
| POST   | `/users`          | Создать пользователя     |
| GET    | `/users`          | Получить список пользователей (`limit`, `offset`) |
| GET    | `/users/{id}`     | Получить пользователя по ID |
| PATCH  | `/users/{id}`     | Частично обновить пользователя |
| DELETE | `/users/{id}`     | Удалить пользователя     |

### Tasks

| Метод  | Путь              | Описание                |
|--------|-------------------|--------------------------|
| POST   | `/tasks`          | Создать задачу           |
| GET    | `/tasks`          | Получить список задач (`limit`, `offset`, `userID`) |
| GET    | `/tasks/{id}`     | Получить задачу по ID    |
| PATCH  | `/tasks/{id}`     | Частично обновить задачу |
| DELETE | `/tasks/{id}`     | Удалить задачу           |

### Statistics

| Метод  | Путь              | Описание                 |
|--------|-------------------|--------------------------|
| GET    | `/statistics`     | Получить статистику (user_id, from, to) |

## Модель данных

- **users**: `id`, `version`, `full_name` (3–100 символов), `phone_number` (опционально, формат `+XXXXXXXXXXX`)
- **tasks**: `id`, `version`, `title` (1–100 символов), `description` (опционально, до 1000 символов), `completed`, `created_at`, `completed_at`, `author_user_id` (ссылка на пользователя)

Схема БД описана в `migrations/000001_init.up.sql`.
