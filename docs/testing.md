# Тестирование

Этот документ описывает, как запускать unit-тесты, интеграционные тесты, нагрузочные сценарии k6 и проверять API через Postman.

- [Unit-тесты](#unit-тесты-usecase-слой)
- [Интеграционные тесты](#интеграционные-тесты)
- [Нагрузочное тестирование](#нагрузочное-тестирование-k6)
- [Проверка сервиса через Postman](#проверка-сервиса-через-postman)


## Unit-тесты (usecase-слой)

Юнит-тесты покрывают usecase-слой сервиса:

- валидация входных DTO (`validators.go`),
- мапперы между доменом и DTO (`mappers.go`),
- методы сервиса:
  - `CreateTeam`, `GetTeam`,
  - `SetIsActive`, `GetUserReviews`,
  - `CreatePullRequest`, `MergePullRequest`, `ReassignReviewer`.

Запуск (через Makefile):

```bash
make unit-test
```

---

## Интеграционные тесты

Интеграционные тесты находятся в `pr-manager-service/internal/integration-tests`.

1. Поднять окружение:

```bash
make dev-up
```

2. Запустить интеграционные тесты:

```bash
make test-integration
```

Тесты проверяют сценарии:

- создание команды и чтение её через `/team/add` + `/team/get`;
- создание PR и получение ревью по пользователю через `/pullRequest/create` + `/users/getReview`.

---

## Нагрузочное тестирование (k6)

Нагрузочные сценарии находятся в `ops/load-testing/`:

- `k6_create_pr.js` — нагрузка на `POST /pullRequest/create`;
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

Подробный отчет по результатам нагрузочного тестирования и сравнению с целевыми SLI находится в файле:

[Отчет по нагрузочному тестированию](load-testing-report.md)

---

## Проверка сервиса через Postman

В каталоге `docs/postman/` лежат:

- коллекция:  
  `pr-manager-service.postman_collection.json`
- окружение:  
  `pr-manager-service-local.postman_environment.json`

Порядок действий:

1. Импортировать окружение `pr-manager-service-local` в Postman.
2. Импортировать коллекцию `pr-manager-service.postman_collection`.
3. Выбрать окружение `pr-manager-service-local`.
4. Поднять backend:

   ```bash
   make dev-up
   ```

5. Выполнить запросы в рекомендуемом порядке:

   - `Create team` (`POST /team/add`)
   - `Get team` (`GET /team/get`)
   - `Set user active/inactive` (`POST /users/setIsActive`)
   - `Create pull request` (`POST /pullRequest/create`)
   - `Reassign reviewer` (`POST /pullRequest/reassign`)
   - `Merge pull request` (`POST /pullRequest/merge`)
   - `Get user reviews` (`GET /users/getReview`)
   - `Stats` (`GET /stats`)
   - `Health` (`GET /health`)

Токены (`admin:<user_id>` / `user:<user_id>`) уже преднастроены в коллекции/окружении.