# PR Manager Service — Avito Test Task (Стажировка Осень 2025)

Сервис для управления pull request’ами и назначениями ревьюеров внутри команд.

Основной стек:
- Go 1.23
- PostgreSQL 16
- Docker / docker compose
- Prometheus + Grafana + Loki + Node Exporter
- k6 для нагрузочного тестирования
- golangci-lint для статического анализа



## Структура проекта

```
.
├── pr-manager-service/
│   ├── cmd/
│   │   └── pr-manager-service/
│   ├── config/
│   ├── internal/
│   │   ├── app/
│   │   ├── adapters/
│   │   │   └── httpadapter/
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── usecase/
│   │   └── integration-tests/
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
│
├── common/
│   └── kit/
│       ├── go.mod
│       └── go.sum
│
├── ops/
│   ├── docker-compose.dev.yml
│   ├── prometheus.yml
│   ├── loki-config.yml
│   ├── promtail-config.yml
│   ├── grafana/
│   │   └── provisioning/
│   │       ├── datasources/
│   │       └── dashboards/
│   └── load-testing/
│       ├── k6_create_pr.js
│       └── k6_get_reviews.js
│
├── docs/
│   ├── contracts/
│   │   └── pr-manager-service-openapi.yml
│   └── postman/
│       ├── PR Reviewer Assignment Service (Test Task, Fall 2025).postman_collection.json
│       └── pr-manager-service-local.postman_environment.json
│
├── .github/
│   └── workflows/
│       └── ci.yml
├── Makefile
├── README.md
└── load-testing-report.md
```
---

## Makefile: основные команды

В корне репозитория находится `Makefile`, который упрощает работу с окружением, тестами и нагрузкой.

### Dev-окружение

```bash
# Запуск dev-окружения (сервис, Postgres, Prometheus, Grafana, Loki, Swagger и т.д.)
make dev-up

# Остановка окружения
make dev-down

# Перезапуск окружения
make dev-restart

# Полная очистка: остановка и удаление volumes
make clear-volumes

# Логи только сервиса
make dev-logs-pr-manager-service

# Логи всех контейнеров docker-compose
make dev-logs-all
```

### Линтер

```bash
# Линтинг сервиса
make lint-pr-manager-service

# Линтинг общего модуля common/kit
make lint-common

# Линтинг всего проекта
make lint
```

`.golangci.yml` лежит в репозитории сервиса, в конфигурации включены проверки:
- форматирование (gofumpt),
- базовые ошибки (errcheck, unused),
- стиль (revive и др.).

### Интеграционные тесты

```bash
# Интеграционные тесты (требуется запущенный сервис: make dev-up)
make test-integration
```

Под капотом команда выполняет:

```bash
cd pr-manager-service && go test ./internal/integration-tests -count=1
```

---

## URL-адреса инфраструктуры

После `make dev-up` доступны:

- Сервис:  
  `http://localhost:8080`

- Swagger UI (документация API):  
  `http://localhost:8082`  
  Использует `docs/contracts/pr-manager-service-openapi.yml`.

- Prometheus:  
  `http://localhost:9090`  
  Можно смотреть сырые метрики, делать запросы PromQL. Сервис экспонирует метрики по `/metrics`.

- Grafana:  
  `http://localhost:3000`  
  Логин/пароль по умолчанию: `admin / admin`.  
  Основные дашборды:
  - **Main Service Metrics** — метрики HTTP-эндпоинтов и бизнес-метрики сервиса.
  - **Node Exporter** — системные метрики хоста (CPU, память, диск).
  - **Logs Dashboard** — дашборд логов из Loki.

- Loki:  
  `http://localhost:3100`  
  Источник логов для Grafana.

- Node Exporter:  
  `http://localhost:9100`  
  Системные метрики, подключены в Prometheus и Grafana.

- Статус/health самого сервиса:
  - `GET /health` — простой healthcheck.
  - `GET /stats` — эндпоинт статистики (имя сервиса, версия, текущее время).

---

## API и аутентификация

Основные эндпоинты сервиса:

- `POST /team/add` — создать команду с участниками.
- `GET  /team/get` — получить команду и список участников.
- `POST /users/setIsActive` — активировать/деактивировать пользователя.
- `GET  /users/getReview` — получить список PR, где пользователь назначен ревьюером.
- `POST /pullRequest/create` — создать PR и автоматически назначить ревьюеров.
- `POST /pullRequest/merge` — пометить PR как смерженный.
- `POST /pullRequest/reassign` — переназначить ревьюера.
- `GET  /stats` — простой эндпоинт статистики сервиса (service name, version, time).
- `GET  /health` — healthcheck.
- `GET  /metrics` — метрики Prometheus.

Аутентификация:

- Заголовок: `Authorization: Bearer <role>:<user_id>`
- Примеры:
  - администратор: `Authorization: Bearer admin:u1`
  - пользователь: `Authorization: Bearer user:u2`

Роль проверяется в HTTP-адаптере (`auth.go`). Для админских операций (создание команд, PR и т.п.) требуется `admin`, для чтения ревью достаточно `user`.

---

## Тестирование

### Unit-тесты (usecase-слой)

Юнит-тесты покрывают usecase-слой в `pr-manager-service/internal/usecase`:

- валидация входных DTO (`validators.go`),
- мапперы между доменом и DTO (`mappers.go`),
- методы сервиса:
  - `CreateTeam`, `GetTeam`,
  - `SetIsActive`, `GetUserReviews`,
  - `CreatePullRequest`, `MergePullRequest`, `ReassignReviewer`.

Запуск:

```bash
cd pr-manager-service
go test ./internal/usecase -count=1
```

`-count=1` отключает кеширование результатов тестов.

### Интеграционные тесты

Интеграционные тесты находятся в `pr-manager-service/internal/integration-tests` и работают поверх запущенного сервиса (по HTTP).

Перед запуском необходимо поднять окружение:

```bash
make dev-up
```

Далее:

```bash
make test-integration
```

Тесты проверяют сценарии:

- создание команды и чтение её через `/team/add` + `/team/get`;
- создание PR и получение ревью по пользователю через `/pullRequest/create` + `/users/getReview`.

---

## Нагрузочное тестирование (k6)

Нагрузочные сценарии находятся в `ops/load-testing/`:

- `k6_create_pr.js` — нагрузка на `POST /pullRequest/create`.
- `k6_get_reviews.js` — нагрузка на `GET /users/getReview`.

Перед запуском нужно поднять сервис:

```bash
make dev-up
```

Затем:

```bash
# сценарий создания PR под нагрузкой
make load-create-pr

# сценарий чтения ревью по пользователю под нагрузкой
make load-get-reviews
```

Под капотом выполняются команды:

```bash
k6 run ops/load-testing/k6_create_pr.js
k6 run ops/load-testing/k6_get_reviews.js
```

Отчёт по результатам нагрузочного тестирования и сравнение с целевыми SLI приведены в отдельном файле `load-testing-report.md`.

---

## Проверка сервиса через Postman

В `docs/` лежат:

- коллекция:  
  `PR Reviewer Assignment Service (Test Task, Fall 2025).postman_collection.json`
- окружение:  
  `pr-manager-service-local.postman_environment.json`

Как использовать:

1. Импортировать окружение `pr-manager-service-local` в Postman.
2. Импортировать коллекцию `PR Reviewer Assignment Service (Test Task, Fall 2025)`.
3. Выбрать окружение `pr-manager-service-local`.
4. Запустить backend:

   ```bash
   make dev-up
   ```

5. Выполнять запросы в рекомендуемом порядке:
   - `Create team` (`POST /team/add`)
   - `Get team` (`GET /team/get`)
   - `Set user active/inactive` (`POST /users/setIsActive`)
   - `Create pull request` (`POST /pullRequest/create`)
   - `Reassign reviewer` (`POST /pullRequest/reassign`)
   - `Merge pull request` (`POST /pullRequest/merge`)
   - `Get user reviews` (`GET /users/getReview`)
   - `Stats` (`GET /stats`)
   - `Health` (`GET /health`)

Токены (`admin:user_id` / `user:user_id`) уже преднастроены в коллекции/окружении.

---

## Stats endpoint

В HTTP-адаптере реализован эндпоинт:

```http
GET /stats
```

Ответ содержит:

```json
{
  "service": "<имя сервиса из cfg.App.Name>",
  "version": "<версия из cfg.App.Version>",
  "time": "<текущее время в формате RFC3339 UTC>"
}
```

Это позволяет быстро проверить:

- имя и версию развернутого сервиса,
- что сервис жив и отвечает по HTTP.

---

## Реализованные требования

Обязательная часть:

- Реализован основной REST API в соответствии с контрактом.
- Хранение данных в PostgreSQL, миграции выполнены через отдельный контейнер `migrate`.
- Разделение слоёв в формате **Чистой архитектуры**: HTTP-адаптер, usecase-слой, хранилище, доменная модель.
- Логирование через общий модуль (`common/kit/logger`).
- Метрики в формате Prometheus, экспорт по `/metrics`.
- Настроены Prometheus, Grafana, Loki, Node Exporter, дашборды для метрик и логов.
- Подготовлена Postman-коллекция и окружение.

Дополнительная часть:

- Добавлен простой эндпоинт статистики: `GET /stats` (имя, версия, время).
- Проведено нагрузочное тестирование решения (k6, два сценария; отчёт в `load-testing-report.md`).
- Реализовано интеграционное/E2E-тестирование (HTTP-тесты в `internal/integration-tests`).
- Описана конфигурация и запуск линтера (`golangci-lint`, gofumpt, errcheck и др.).

Массовая деактивация пользователей команды и расширенная безопасная переназначаемость открытых PR не реализованы в рамках тестового, но архитектура сервисов и usecase-слоя позволяет добавить эту функциональность поверх существующих интерфейсов.