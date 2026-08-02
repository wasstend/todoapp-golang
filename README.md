# To Do App

Простой REST API для управления задачами, написанный на Go, с простым веб-интерфейсом. Приложение построено по принципу "feature-based" архитектуры: каждая доменная область (`users`, `tasks`, `statistics`, `web`) содержит свои слои `service` / `repository` / `transport`, а общий код (логгер, HTTP-сервер, работа с Postgres) вынесен в `internal/core`.

## Стек технологий

- **Go** 1.26
- **PostgreSQL 18** — хранилище данных, доступ через [pgx/v5](https://github.com/jackc/pgx)
- **golang-migrate** — миграции БД
- **zap** — структурированное логирование
- **go-playground/validator** — валидация входных данных
- **net/http** (`http.ServeMux`) — HTTP-роутинг, без сторонних фреймворков
- **swaggo/swag** + **http-swagger** — генерация и раздача Swagger-документации
- **Docker Compose** — окружение для Postgres, миграций и деплоя приложения

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
  ├── statistics/      — статистика
  └── web/             — раздача статичной веб-страницы (public/index.html)
docs/                  — сгенерированная Swagger-документация (swag)
public/                — статика веб-интерфейса
migrations/            — SQL-миграции (golang-migrate)
docker-compose.yaml    — Postgres, туннель для локального порта, миграции, сборка/деплой приложения, генерация Swagger
Makefile               — команды для запуска окружения, миграций, деплоя и приложения
```

## Настройка окружения

1. Скопируйте файл с переменными окружения и заполните его:

   ```bash
   cp .env.example .env
   ```

   Основные переменные:

   | Переменная              | Описание                                         | По умолчанию |
   |-------------------------|---------------------------------------------------|--------------|
   | `POSTGRES_USER`         | пользователь Postgres                              | —            |
   | `POSTGRES_PASSWORD`     | пароль Postgres                                    | —            |
   | `POSTGRES_DB`           | имя базы данных                                    | —            |
   | `POSTGRES_TIMEOUT`      | таймаут подключения к БД                           | `10s`        |
   | `HTTP_ADDR`             | адрес, на котором слушает HTTP-сервер              | `:5050`      |
   | `HTTP_SHUTDOWN_TIMEOUT` | таймаут graceful shutdown сервера                  | `30s`        |
   | `HTTP_ALLOWED_ORIGINS`  | список разрешённых CORS origin'ов через запятую    | —            |
   | `LOGGER_LEVEL`          | уровень логирования                                | `DEBUG`      |

## Запуск

### Простой запуск

Прострой запуск доступен по **URL:**`http://5.101.50.146:5050`

### Поднять локально

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

   По умолчанию сервер поднимется на `http://localhost:5050`, логи пишутся в `out/logs/`. Веб-интерфейс доступен на `/`, Swagger-документация — на `/swagger/`.

   Приложение можно также запустить в Docker вместо `make todoapp-run`:

   ```bash
   make todoapp-deploy    # собрать образ и поднять контейнер todoapp
   make todoapp-undeploy  # остановить контейнер todoapp
   ```

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

### Swagger

Документация генерируется из аннотаций в коде (`swaggo/swag`) в директорию `docs/`:

```bash
make swagger-gen
```

После запуска приложения документация доступна на `http://localhost:5050/swagger/`.

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

Вне префикса `/api/v1` также доступны:

| Метод  | Путь              | Описание                 |
|--------|-------------------|--------------------------|
| GET    | `/`               | Веб-интерфейс (`public/index.html`) |
| GET    | `/swagger/`       | Swagger UI               |

## Модель данных

Схема БД описана в `migrations/000001_init.up.sql`.
