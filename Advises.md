# Сборник: советы и ошибки при создании микросервиса на Go

## text/Solution and advice (1).md
_Размер: 14478 байт_
_Кодировка: utf-8_

```
# Стажировки для Backend-разработчиков в Авито: на что мы смотрим при проверке работ

Наши эксперты проанализировали тестовые задания кандидатов прошедшего отбора и рассказали о том, какие ошибки встречались чаще всего.  
Все совпадения случайны!

---

## Валидация бизнес-правил

Во многих проектах бизнес-правила — например, допустимые города или типы товаров — проверялись только на уровне базы данных.  

**Пример:**

```sql
city VARCHAR(100) NOT NULL
   CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
type VARCHAR(50) NOT NULL
   CHECK (type IN ('electronics', 'clothing', 'shoes')),
status VARCHAR(20) NOT NULL
   CHECK (status IN ('in_progress', 'closed'))
   DEFAULT 'in_progress'
```
На первый взгляд кажется, что решение простое и надёжное: база не позволит записать некорректные данные. Однако на практике это создаёт больше проблем, чем пользы.

---

### Почему так не стоит делать

1. **Нарушается принцип разделения ответственности**  
   Бизнес-правила — это часть прикладной логики, а не уровня хранения данных. Когда их вшивают в схему, валидация расползается между слоями.

2. **Становится неудобно сопровождать базу данных**  
   Чтобы изменить бизнес-правила, приходится создавать миграции, перекатывать схему, пересоздавать constraint. Особенно сложно, когда бизнес-правила часто меняются — к примеру, добавляется много новых городов или типов товаров.

3. **Некорректно обрабатываются ошибки**  

**Какая ошибка должна быть:**  
```json
{"error": "город 'Томск' недоступен для открытия ПВЗ"}
```

**Какой она будет, если вшить бизнес-правила в базу данных:**  
```
pq: new row for relation "pvz" violates check constraint "allowed_city"
```

Обработку ошибок можно исправить дублированием логики в приложении — но в таком случае теряется единый источник истины и повышается риск расхождений.

4. **Такой код сложнее тестировать**  
   Если валидация проведена только на уровне базы данных, её невозможно проверить в отрыве от инфраструктуры.

5. **Базу данных становится сложно переносить**  
   СУБД по-разному реализуют ограничения, а в некоторых их и вовсе может не быть.

---

### Как тогда делать

Рекомендуем проводить валидацию на уровне приложения. База данных должна обеспечивать только целостность — например, `NOT NULL` и внешние ключи.  

В редких случаях в базе может быть дополнительная защита — но только если вы уверены, что список будет меняться редко. Даже тогда это не должно быть единственным местом проверки.

---

## Хранение данных

В некоторых проектах поля `city`, `type` и `status` хранились как `VARCHAR` или `TEXT`. На первый взгляд кажется, что это просто и удобно, но у такого подхода есть серьёзные недостатки.

---

### Почему так не стоит делать

1. **Становится сложно вносить изменения**  
   Если понадобится переименовать значение — например, `"clothing"` в `"clothing/t-shirt"` — придётся обновлять все записи в таблице. Это долго, ресурсоёмко и повышает риск рассинхронизации.

2. **Снижается эффективность кода**  
   Строки занимают больше места, чем числовые идентификаторы. Операции сравнения строк выполняются медленнее, чем сравнение чисел, особенно на больших объёмах данных.

3. **Появляются трудности в индексации и агрегации**  
   Индексы по строковым полям больше по размеру и работают медленнее. Это снижает производительность фильтрации, сортировки и агрегирующих запросов.

---

## Почему не стоит полагаться на ENUM

Другой распространённый вариант — использовать `ENUM` вместо строк. У такого подхода есть плюсы, но и существенные минусы.

1. **Чтобы изменить или удалить значение, придётся перестраивать всю таблицу.**  
   Если таблица большая, это приводит к блокировкам и простоям.

2. **Ограничивается гибкость.**  
   При добавлении новых значений приходится постоянно обновлять миграции.

---

### Как тогда делать

Рекомендуем провести нормализацию данных. Самое устойчивое решение — вынести справочники (`city`, `type`, `status`) в отдельные таблицы:  

- Таблица **cities** (`id`, `name`)  
- Таблица **product_types** (`id`, `name`)  
- Таблица **statuses** (`id`, `name`)  

А в основной таблице хранить только целочисленные ключи: `city_id`, `type_id`, `status_id`.

**Плюсы подхода:**  
- Быстрое хранение и сравнение.  
- Возможность централизованно менять значения.  
- Удобство индексации и агрегации.  
- Гибкость структуры.  

---

## Работа с индексами

Во многих решениях кандидаты не добавляли индексы к таблицам. На тестовых данных это может быть незаметно, но в реальной системе быстро приводит к проблемам.

---

### Почему так не стоит делать

1. **Снижается производительность запросов.**  
   При фильтрации или `JOIN` без индексов база сканирует всю таблицу.  

2. **Усложняется масштабирование.**  
   При росте объёма данных запросы начинают выполняться значительно дольше.

---

### Как тогда делать

В таблицах важно добавлять индексы для:  
- Полей, используемых в фильтрах (`WHERE city_id = ? AND date BETWEEN …`),  
- Внешних ключей и связей между таблицами,  
- Полей сортировки и пагинации (`ORDER BY created_at, id`).  

Для сложных фильтров — используйте составные индексы, например `(pvz_id, date)`.

⚠️ **Важно:** избыточные индексы замедляют записи (`INSERT`, `UPDATE`, `DELETE`) и занимают место. Добавляйте их только после анализа `EXPLAIN ANALYZE`.

---

## Тестирование кода

Во многих работах код не был протестирован. Часто кандидаты покрывают тестами инфраструктурные части, но не бизнес-логику.

---

### Как делать не надо

```go
func TestGetUserRepository(t *testing.T) {
   repo := NewUserRepo(mockDB)
   _, err := repo.GetByID(42)
   require.NoError(t, err)
}
```

---

### Почему так не стоит делать

1. **Создаётся ложное чувство защищённости.**  
   Такие тесты не ловят ошибки в логике.  

2. **Смещается фокус внимания.**  
   Проверяются детали реализации, а не сценарии.  

3. **Повышаются трудозатраты.**  
   Инфраструктурные тесты часто ломаются при рефакторинге.  

---

### Как тогда делать

Сосредоточьтесь на тестировании бизнес-логики:  
- Переходы состояний (например, статусы заказов).  
- Расчёты скидок и комиссий.  
- Обработку ошибок и граничных условий.  

---

## «Толщина» интерфейса

Иногда встречаются интерфейсы, объединяющие слишком много сущностей.  

---

### Как делать не надо

```go
type Repository interface {
    CreateUser(...)
    GetUserByEmail(...)
    CreatePVZ(...)
    GetReceptionByID(...)
    CreateProduct(...)
}
```

---

### Как тогда делать

Разделяйте интерфейсы по доменам — `UserRepository`, `PVZRepository`, `ProductRepository`.  
Это повышает читаемость, тестопригодность и гибкость архитектуры.

---

## Собранность бизнес-логики

Бизнес-логика должна быть сосредоточена в одном сервисном слое. Когда проверки и бизнес-правила размазаны по обработчикам, репозиториям и вспомогательным функциям — код становится нечитабельным.

---

### Как тогда делать

- Бизнес-правила — в сервисах.  
- Репозитории — только за доступ к данным.  
- Handlers — только за приём и отдачу запросов.

---

## Лаконичность в транзакциях

Не стоит оборачивать каждую операцию в транзакцию.  

---

### Как делать не надо

```go
func (r *PostgresRepo) GetLastOrder(ctx context.Context, warehouseID int) (*Order, error) {
   var order Order
   err := r.withTx(ctx, func(tx *sqlx.Tx) error {
       return tx.QueryRowContext(
           ctx,
           `SELECT id, status FROM orders WHERE warehouse_id = $1 ORDER BY created_at DESC LIMIT 1`,
           warehouseID,
       ).Scan(&order.ID, &order.Status)
   })
   if err != nil {
       return nil, fmt.Errorf("failed to get last order: %w", err)
   }
   return &order, nil
}
```

---

### Как тогда делать

Используйте транзакции только там, где нужно сохранить атомарность нескольких изменений. Решение о применении транзакции должно приниматься на уровне бизнес-логики.

---

## Работа с конфигурацией

Нельзя хранить чувствительные данные в коде.

---

### Как делать не надо

```go
sql.Open("postgres", "postgres://user:password@localhost:5432/master")
```

---

### Как тогда делать

Используйте переменные окружения:

```go
sql.Open("postgres", os.Getenv("DB_CONN"))
```

**Примеры:**  
- `export DB_CONN="postgres://user:password@localhost:5432/master"`  
- `.env` файл  
- Docker/Kubernetes параметры  
- Vault или другое хранилище секретов

---

## Настройка линтеров

Несоблюдение форматирования и правил стиля снижает читаемость и качество кода.

---

### Как делать не надо

```go
func add(a int, b int) int {
   return a+b
}
```

**При запуске линтера:**  
```
gofmt: File is not formatted with gofmt
```

---

### Как тогда делать

Настройте линтеры и автоформатирование (например, `pre-commit hook`).  
Регулярно обновляйте конфигурацию линтеров.

---

## Логирование в приложении

Ошибки нужно логировать, а не игнорировать.

---

### Как делать не надо

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
   if _, err := db.Query("INSERT INTO users(name) VALUES('test')"); err != nil {
       w.WriteHeader(http.StatusBadRequest)
   } else {
       w.WriteHeader(http.StatusOK)
   }
}
```

---

### Как тогда делать

Используйте логирование (например, Zap или Slog).  
Не логируйте чувствительные данные (пароли, токены).

---


```

## text/Solution and advice (2).md
_Размер: 13690 байт_
_Кодировка: utf-8_

```
# Стажировки для Backend-разработчиков в Авито: на что мы смотрим при проверке работ — первая часть

Наши эксперты отсмотрели много тестовых заданий кандидатов — и теперь делятся распространёнными ошибками, чтобы вы могли учесть их в будущем.

## 1. Code style и структура кода

### 1.1. Сортировка импорта
В некоторых проектах отсутствовал единый автоформат, были перемешаны импорты, а отступы и переносы строк — не согласованы.

**Почему так не стоит делать**  
В таком случае падает читаемость кода и скорость код-ревью. Растёт вероятность случайных конфликтов и мелких ошибок при правках.

**Как правильно**  
Стоит включить автоформат и линтер.

---

### 1.2. Разделение обработчиков (handlers)
В работах многих кандидатов обработчики были собраны в одном файле. В `handlers.go` лежали все эндпоинты — `/api/auth`, `/api/info` — без логического разделения по доменам.

**Почему так не стоит делать**  
В такой ситуации файл разрастается, усложняется навигация, растёт риск конфликтов и утечек ответственности между обработчиками.

**Как правильно**  
Стоит разложить эндпоинты по ресурсам или доменам.

**Хороший пример:**
```bash
/internal/http/handlers/
  auth.go       // /api/auth
  coins.go      // /api/sendCoin
  ...
```

---

### 1.3. Моки и код в проде
Часто моки лежат там же, где основной код.

**Почему так не стоит делать**  
В таком случае в навигации появляется шум, растёт риск того, что в проде случайно используется мок, а зависимости пересекутся.

**Как правильно**  
Хранить моки стоит в `mocks/`, а генерировать — через `mockery`. Больше деталей ищите в документации.

---

### 1.4. Названия структур внутри функций
Иногда в функциях создаются «сырые» `struct{...}` вместо именованных типов.

**Почему так не стоит делать**  
Не получается переиспользовать типы, трудно проводить валидацию и тестирование.

**Как правильно**
Стоит ввести именованные DTO или domain-модели.

**Плохой пример**
```go
req := struct{ Name string }{Name: "John"}
```

**Хороший пример**
```go
type User struct {
  Name string
}
```

---

## 2. Архитектура и разделение слоёв

### 2.1. Бизнес-логика и слои
В некоторых проектах HTTP-слой содержал бизнес-логику и выполнял операции с базами данных.

**Почему так не стоит делать**  
В такой ситуации код трудно тестировать, повторно использовать и менять логику без изменений в HTTP-handlers.

**Как правильно**  
Стоит оставлять валидацию входных данных на уровне `handlers`, а бизнес-логику реализовывать в `usecase` — то есть придерживаться чистой архитектуры.

**Хороший пример:**
```go
// Handler: только парсинг и маппинг.
func (h *Handler) SendCoin(w http.ResponseWriter, r *http.Request) {
    var in SendCoinInput
    ctx := r.Context()

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "bad request", 400)
        return
    }

    if err := h.usecase.SendCoin(ctx, in.From, in.To, in.Amount); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}
```

---

### 2.2. Зависимость билдера
Часто бывает так, что билдер импортирует реализацию, а не интерфейс.

**Почему так не стоит делать**  
Появляется жёсткая связность. В таком случае не получится мокать зависимость в тестах и менять реализацию.

**Как правильно**  
Стоит объявлять интерфейсы на стороне потребителя и создавать зависимость от интерфейса.

**Хороший пример:**
```bash
/internal/repo/merch/
  merch.go                // Реализация репозитория merch.
/internal/service/merch/
  merch.go                // Реализация сервиса merch.
```

```go
// /internal/service/merch/merch.go
package merch

type MerchRepo interface {
    Get(ctx context.Context, name string) error
}

type MerchService struct {
    Repo MerchRepo
}
```

---

### 2.3. Миграции
Иногда кандидаты пишут миграции баз данных полностью в `init.sql`.

**Почему так не стоит делать**  
Становится сложно откатывать или перекатывать миграции частями, теряется история изменений.

**Как правильно**  
Стоит создавать отдельные миграции, а для управления ими использовать специальный инструмент — например, golang-migrate.

**Хороший пример:**
```bash
/migrations/
  0001_init.sql
  0002_add_merch.sql
```

---

## 3. Тесты

### 3.1. Юнит-тесты и e2e
Проверяйте, что код покрыт тестами и они не падают, не появляются flaky тесты.

**Как правильно**  
Стоит добиваться зелёных тестов и исправлять причины flaky тестов. Например, они могут появляться при проверке равенства map — порядок ключей всегда разный.

---

### 3.2. Отдельное тестовое окружение
Следите за тем, чтобы тесты не использовали ту же базу данных, что и приложение. Иначе данные сломаются, и тесты станут нестабильными.

**Как правильно**  
Стоит поднимать отдельную тестовую базу данных. А e2e-тесты лучше выделить в отдельный `docker-compose.e2e.yaml` со своими базами данных, переменными и сетью.

**Хороший пример:**
```yaml
# docker-compose.e2e.yaml
services:
  db_e2e:
    image: postgres:15
    environment:
      POSTGRES_DB: database
      POSTGRES_PASSWORD: password
  api_e2e:
    build: .
    env_file: .env.e2e
    depends_on: [db_e2e]
  tests:
    build: ./e2e_tests
    depends_on: [api_e2e]
```

---

## 4. DevOps и окружение

### 4.1. Работа `docker compose up`
В некоторых проектах сервис не поднимается без ручных правок, в CI не запускается окружение.

**Как правильно**  
Стоит регулярно проверять чистый старт, добавлять healthchecks и зависимости, прописывать entrypoint для миграций.

---

### 4.2. Совпадение версий Go в `go.mod` и Dockerfile
Бывает, в `go.mod` указана одна версия Go, а образ сборки в Dockerfile — другой. В таком случае сборка падает.

**Как правильно**  
Стоит синхронизировать версии и использовать multi-stage build.

**Хороший пример:**
```dockerfile
FROM golang:1.22

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /build ./internal/cmd     && go clean -cache -modcache

EXPOSE 8080
CMD ["/build"]
```

---

### 4.3. Заголовки для JWT
Иногда кандидаты кладут токен в произвольный заголовок вместо `Authorization: Bearer <token>`. Это приводит к несовместимости с библиотеками и путанице в клиентах и тестах.

**Как правильно**  
Стоит придерживаться стандартных заголовков.

**Плохой пример:**
```http
GET /api/info HTTP/1.1
Host: localhost:8080
X-Auth: eyJhbGciOiJIUzI1NiIsInR5cCI6...
```

**Хороший пример:**
```http
GET /api/info HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6...
```

```go
// Чтение токена.
auth := r.Header.Get("Authorization")

const prefix = "Bearer "

if !strings.HasPrefix(auth, prefix) {
  http.Error(w, "unauthorized", http.StatusUnauthorized)
  return
}

token := strings.TrimPrefix(auth, prefix)
// Далее валидация token...
```

---

## 5. Безопасность

### 5.1. Коммиты `.env` в репозиторий
Ошибка, которая встречается очень часто: секреты — к примеру, JWT-ключ или пароли — лежат в репозитории. Это может приводить к серьёзным утечкам данных.

**Как правильно**  
Стоит добавлять `.env` в `.gitignore`, хранить только `.env.example`. Также важно использовать переменные окружения или секрет-менеджер.

**Хороший пример:**
```go
type Config struct {
    JWTKey string `env:"JWT_KEY,required"`
}
```

---

### 5.2. Хеширование паролей
Ещё одна частая ошибка, когда пароли пользователей хранятся в базе данных в открытом виде. Это может приводить к компрометации и репутационным рискам.

**Как правильно**  
Стоит хранить bcrypt-хеши.

---

## 6. Баги в бизнес-логике

### 6.1. Накрутка монет
В некоторых проектах пользователь может отправить монеты сам себе: такие лазейки обычно приводят к накруткам баланса и фроду.

**Как правильно**  
Стоит добавлять валидацию `from user != to user`.

**Хороший пример:**
```go
if fromUserID == toUserID { return ErrSelfTransfer }
```

---

### 6.2. Неатомарные операции
В некоторых проектах баланс может меняться без транзакции.

**Как правильно**  
Стоит оборачивать бизнес-операции в транзакции базы данных. Для этого нужно использовать transaction-manager.

**Хороший пример:**
```go
err := trManager.Do(ctx, func(ctx context.Context) error {
    balance, err := GetBalance(ctx, fromUser)
    if err != nil {
        return err
    }

    if balance < amount {
        return ErrInsufficientFunds
    }

    if err := AddBalance(ctx, fromUser, -amount); err != nil {
        return err
    }

    if err := AddBalance(ctx, toUser, amount); err != nil {
        return err
    }

    return nil
})
```

---

### 6.3. Сравнивание ошибок
Бывает, ошибки сопоставляются через `==` или `strings.Contains()`. В таком случае растёт риск пропустить специфичную ошибку — а проверить обёрнутую ошибку и вовсе становится невозможно.

**Как правильно**  
Стоит использовать `errors.Is()`.

**Хороший пример:**
```go
if errors.Is(err, ErrInsufficientFunds) { /* ... */ }
```

---

## 7. Технический долг

### 7.1. `TODO` без привязки к задаче
Оставлять комментарии-заглушки без номера тикета и срока — это моветон. Так статус кода становится неясным, копятся долги.

**Как правильно**  
Стоит переносить todo в issue-трекер, а в коде ссылаться на задачу.

**Хороший пример:**
```go
// TODO (TASK-123): Заменить поход в БД на кэш к 01.01.26.
```

```

## text/Solution and advice.md
_Размер: 20695 байт_
_Кодировка: utf-8_

```
# Стажировки для Backend-разработчиков в Авито: на что мы смотрим при проверке работ
Наши разработчики проанализировали тестовые задания кандидатов прошедшего отбора и написали статью о лучших подходах и наиболее частых ошибках в работах.  Все совпадения, конечно же, случайны :)
# Инструкция для запуска проекта
Проверяющие должны просмотреть много кода и проверить работоспособность нескольких проектов.  
Помогите им быстрее перейти к интересному — к вашему коду и тестированию логики приложения. Для этого подготовьте понятную инструкцию по проекту.  
Это точно будет отмечено при проверке и поможет вам выделиться среди других кандидатов.
## О чём можно написать
Как запустить проект. Например:  
```markdown
Для запуска проекта нужно выполнить команду `docker-compose up`.
После этого сервис будет доступен на порту `:8080`
```
## Проверка
Перед отправкой нужно проверить, что проект без проблем запускается с 0:  
- указанный в инструкции скрипт завершается без ошибок;
- миграции базы данных успешно проходят;  
- сервис запускается;  
- API сервиса работает.
# Проект должен запускаться
Проекты некоторых кандидатов не запускались на компьютере проверяющего.
Проверьте, что ваш сервис запускается «с нуля», а базу данных легко мигрировать.  
Правилом хорошего тона считается добавить Dockerfile и параметры окружения (ENV-параметры), чтобы сервис можно было запускать на любой системе.
# Формат ответа
Часто кандидаты забывают соблюсти формат ответа, который указан [в задании](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-autumn-2024/openapi.yml).
В итоге возвращается ответ в формате `json`, даже когда нужен `text/plain`,  или для некоторых случаев не проставляется нужный заголовок.
Исправить это просто — внимательность и  
```go  
w.Header().Set("Content-Type", "application/json")  
```
# Валидация входных данных
Если не проверять входные данные на корректность,   пользователь может случайно (или намеренно) повредить или получить доступ к данным других юзеров,  а в худшем случае — ко всей базе данных.
## Базовые валидации
 Добавляйте валидации в контроллерах при парсинге запросов, чтобы  по системе передавались  корректные данные.  
Если валидация не прошла — сразу возвращать ошибку.
Например:  
```go

type TenderServiceType string

const (  
	TenderServiceTypeConstruction TenderServiceType = "Construction"  
	TenderServiceTypeDelivery     TenderServiceType = "Delivery"  
	TenderServiceTypeManufacture  TenderServiceType = "Manufacture"  
)

func (tst TenderServiceType) Validate() bool {  
	switch tst {  
	case TenderServiceTypeConstruction, TenderServiceTypeDelivery, TenderServiceTypeManufacture:  
		return true  
	default:  
		return false  
	}  
}

// ...

types := r.URL.Query()["service_type"]

serviceTypes := make([]TenderServiceType{}, 0, len(types))  
for _, t := range types {  
	tst := TenderServiceType(t)  
	if !tst.Validate() {  
		continue  
	}  
	serviceTypes = append(serviceTypes, tst)  
}

if len(serviceTypes) == 0 {  
	// возвращаем ошибку  
}  
```

## Специфические валидации
Важно не забывать про отношения между сущностями.  
Допустим, пользователь Вася авторизован, но у него нет прав на изменение предложения, которое создал пользователь Петя.
**Для таких проверок можно или делать дополнительный запрос в базу данных**:
```go  
func (s *Storage) checkRelationToOrganization(ctx context.Context, userId, orgId uuid.UUID) bool {  
	res := 0  
	query := `SELECT 1 FROM organization_responsible WHERE user_id = $1 AND organization_id = $2;`  
	_ = s.conn.QueryRow(ctx, query, userId, orgId).Scan(&res) // ошибку нужно обязательно обработать  
	return res > 0  
}  
```

или при обычном запросе добавить дополнительное условие.
** о построении запросов к базе данных есть отдельный раздел.
Особенно важно проверять авторизацию пользователя.  Подробнее можно почитать в [разборе заданий предыдущей волны](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2024/Solution%20and%20advice/%D0%90%D0%B2%D1%82%D0%BE%D1%80%D0%B8%D0%B7%D0%B0%D1%86%D0%B8%D1%8F.md)
# Code Style
Следование общепринятым рекомендациям по стилю кода делает его более понятным, читаемым, поддерживаемым, упрощает внесение изменений.

Наиболее частые проблемы:  
- не используется линтер, например golangci-lint, который мог бы исправить часть допущенных ошибок;  
- аналогично для языка python не соблюдается PEP8, не используется линтер, например flake8;  
- магические строки (не вынесены в константы);  
- нет логических разделений между блоками;  
- мега-структуры/интерфейсы с большой зоной ответственности;  
- не оборачиваются ошибки;  
- импорты не отсортированы;  
- не пробрасывается контекст.
# Меньше повторяющегося кода
Проверяющий смотрит не только, насколько правильно работает код, но и как легко его поддерживать в дальнейшем.
Если в коде часто встречаются дубли (много одинаковых строк), такой код сложнее поддерживать, он больше подвержен ошибкам.  
Старайтесь выносить повторяющуюся логику в отдельные функции.
Например:

```go  
func respondWithError(w http.ResponseWriter, statusCode int, handlerName string, err error) {  
	slog.Error(err.Error(), "handler", handlerName)

	w.Header().Set("Content-Type", "application/json")  
	w.WriteHeader(statusCode)  
	_, _ = w.Write([]byte(`{"error":"Error occurred"}`))  
}  
```
# Бизнес-логика на своём месте
Часто кандидаты помещают всю логику или на уровень хэндлера (контроллера), или на уровень работы с базой данных.
Это лишает приложение гибкости. Лучше разделять его на разные слои.  
- Если транспортный слой (хэндлеры) содержит основную логику, то её сложно переиспользовать,если в дальнейшем мы захотим поменять или добавить новый протокол (например, grpc+rest api).  
- Если логика вынесена на уровень инфраструктуры (баз данных), то это затруднит работу с разными хранилищами.  
Такую логику обычно выносят на уровень сервисов.
Следует ограничивать взаимодействие между слоями, например, использование базы данных в хендлере.  
Для управления направлением зависимостей между слоями лучше использовать интерфейсы.

Хороший пример:  
```go  
// internal/usecase/tender/deps.go  
type TenderRepository interface {  
    GetTenderById(ctx context.Context, tenderID uuid.UUID) (model.Tender, error)  
}

// internal/usecase/tender.go  
type TenderUsecase struct {  
    tenderRepository TenderRepository  
}

// internal/repository/tender/repository.go  
func (r *Repository) GetTenderById(ctx context.Context, id uuid.UUID) (tender.Tender, error) {  
    // implementation of TenderRepository  
}  
```

Интерфейс `TenderRepository` можно использовать для генерации моков для тестирования usecase-а.

На тему разделения логики на слои можно почитать про Чистую или Гексагональную Архитектуры.

Также рекомендуем почитать [совет для прошлой волны](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2024/Solution%20and%20advice/%D0%91%D0%B8%D0%B7%D0%BD%D0%B5%D1%81-%D0%BB%D0%BE%D0%B3%D0%B8%D0%BA%D0%B0%20%D0%B2%20handler.md).

# Код покрыт тестами

Самый распространённый недочёт в этой волне — отсутствие тестов.

Задание было очень объёмным,но всё-таки стоитнаписать хотя бы несколько тестов для того, чтобы показать свои навыки и понимание темы.  
В идеале покрыть тестами несколько типичных методов: хэндлер, сервис, метод для работы с базой данных.

Лучше использовать табличные тесты, это обычно позволяет сократить количество кода в тестовых файлах.

# Логирование ошибок

В реальных системах, использующих микросервисы, применяются 4 вида сигналов для наблюдения за состоянием системы:логи, метрики, алерты и распределённый трейсинг.

Проще всего добавить в свой сервислоги.
Сделать это можно, например, так:
Инициализируем логгер в файле `main.go` и передаём его в конструкторы, например, сервисов.  
```go  
lgr := slog.New(slog.NewJSONHandler(os.Stderr, nil))

serviceBid := service.NewBidService(lgr)  
```

При обработке запросов внутри сервисов мы сможем логировать важную информацию.

```go  
func NewBidService(lgr *slog.Logger) *BidService {  
	return &BidService{  
		lgr: lgr,  
	}  
}

func (s *BidService) CreateBid(ctx context.Context, data *Request) (Response, error) {  
	if bid.Name == "" {  
		s.lgr.With(  
			slog.String("username", data.UserName),  
		).Error("creating a bid: username is empty")

		return Response{}, ErrUsernameFieldEmpty  
	}  
}  
```

Конечно, в сервисах лучше использовать интерфейсы на сущности, это упрощает тестирование через мок-объекты.  
Но это уже немного другая тема.

# Пользователь не должен видеть полные ошибки

Часто пользователю отправляются ошибки напрямую из, например, базы данных.  Ему такая информация не нужна.  
Более того, таким образом мы облегчаем злоумышленнику задачу взлома нашей базы данных - ему будет легче сразу видеть результат своих действий.  
Особенно это опасно, если мы не озаботились валидацией пользовательских данных и правильной генерацией SQL-запросов.

```go  
func (s *BidService) CreateBid(ctx context.Context, data *Request) (Response, error) {  
	err := s.db.CreateBid(ctx, data)  
	if err != nil {  
		s.lgr.With(  
			slog.Any("username", bid.UserName),  
		).Error("creating a bid: " + err.Error())

		return Response{}, InternalError  
	}  
}  
```
Выше показан упрощённый пример. Обычно ошибки подменяются на пользовательские на транспортном уровне.

# Грамотное использование транзакций

При добавлении или изменении данных в базе важно использовать транзакции. Часто изменения касаются нескольких таблиц в рамках одной операции со стороны пользователя. Чтобы сохранить консистентность данных, все изменения лучше делать в рамках одной транзакции.

Пример того, как это можно реализовать:  
```go  
func (s *Storage) SubmitBidDecision(ctx context.Context, bidID, tenderID uuid.UUID) (err error) {  
	tx, err := s.conn.Begin(ctx)  
	if err != nil {  
		return fmt.Errorf("starting transaction: %w", err)  
	}

	defer func() {  
		var e error  
		if err == nil {  
			e = tx.Commit(ctx)  
		} else {  
			e = tx.Rollback(ctx)  
		}

		if err == nil && e != nil {  
			err = fmt.Errorf("finishing transaction: %w", e)  
		}  
	}()

	queryBid := `  
UPDATE bid   
SET   
	status = $1  
WHERE id = $2 ;`

	if _, err = tx.Exec(ctx, queryBid, model.BidStatusApproved, bidID); err != nil {  
		return fmt.Errorf("updating bids: %w", err)  
	}

	queryTender := `  
UPDATE tender   
SET   
	status = $1  
WHERE id = $2;`

	if _, err = tx.Exec(ctx, queryTender, model.TenderStatusClosed, tenderID); err != nil {  
		return fmt.Errorf("updating tenders: %w", err)  
	}

	return nil  
}  
```

*вопросу правильного построения запросов к базе данных посвящён отдельный совет.

При этом не стоит использовать транзакции всегда.  
Об этом можно прочитать в [совете для предыдущей волны](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2024/Solution%20and%20advice/%D0%9F%D0%BE%D0%B2%D1%81%D0%B5%D0%BC%D0%B5%D1%81%D1%82%D0%BD%D0%BE%D0%B5%20%D0%B8%D1%81%D0%BF%D0%BE%D0%BB%D1%8C%D0%B7%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D0%B5%20%D1%82%D1%80%D0%B0%D0%BD%D0%B7%D0%B0%D0%BA%D1%86%D0%B8%D0%B9.md)

## Альтернативные варианты работы с транзакцией

Кроме ручной работы с транзакциями в слое базы данных можно использовать менеджер транзакций,у нас есть отличная [статья на хабре про него](https://habr.com/ru/companies/avito/articles/727168/).

## Защита от ошибок при конкурентном использовании (продвинутый кейс)

Пример выше не гарантирует, что конкурентная транзакция уже не изменила статус тендера.
Чтобы обеспечить корректную работу сервиса в этой ситуации, необходимо дополнительно проработать логику:

- использовать явные блокировки https://postgrespro.ru/docs/postgrespro/17/explicit-locking;  
- полагаться на уровни изоляции в транзакции, например, создать транзакцию с уровнем Repeatable read/Serializable,внутри транзакции запросить статус тендера и рассчитывать на то,  
что база данных автоматически отклонит транзакцию в случае, если параллельная транзакция изменит статус;  
- в рамках одного запроса сделать проверку статуса и его обновление, с этим может помочь CTE.

# Построение запросов к базе данных

Для безопасной вставки пользовательских данных в базу нужно их экранировать.  Обычно для простоты экранируют все данные, которые добавляются/изменяются в базе.

Также для переиспользования одного query в разных SQL-запросах (например, с разными условиями `where`, для использования в подзапросах и т.п.),  
часто прибегают к помощи SQL generators или SQL Query Builders.

В рамках тестового задания использование генераторов может быть чрезмерным, поэтому проще воспользоваться билдером запросов, например, [Squirrel](https://github.com/Masterminds/squirrel).

Пример использования:  
```go  
func (s *Storage) SubmitBidDecision(ctx context.Context, bidID, tenderID uuid.UUID) (err error) {  
	// ...

	// достаточно инициализировать builder один раз и присвоить структуре.  
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.Update("bid").  
		Set("status", model.BidStatusApproved).  
		Where(squirrel.Eq{"id": bidID}).  
		ToSql()  
	if err != nil {  
		return fmt.Errorf("building query: %w", err)  
	}

	_, err = s.conn.Exec(ctx, query, args...)  
	if err != nil {  
		return fmt.Errorf("executing query: %w", err)  
	}

	// ...

	return nil  
}  
```

Таким образом мы гарантируем, что данные для вставки будут экранированы. Более того, в Squirrel есть приятные бонусы, например, кэширование prepared statements.

Эта тема также обсуждалась в [советах к прошлой волне](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2024/Solution%20and%20advice/%D0%9F%D0%BE%D1%81%D1%82%D1%80%D0%BE%D0%B5%D0%BD%D0%B8%D0%B5%20%D0%B4%D0%B8%D0%BD%D0%B0%D0%BC%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D1%85%20%D0%B7%D0%B0%D0%BF%D1%80%D0%BE%D1%81%D0%BE%D0%B2%20%D0%BA%20%D0%91%D0%94.md).  

```

## text/Авторизация.md
_Размер: 2258 байт_
_Кодировка: utf-8_

```
# Авторизация

**Ошибка:** 

Отсутствие авторизации или использование одного токена доступа для всех пользователей, пример:

```go
const (
	adminToken = "admin_token"
	userToken  = "user_token"
)

func authenticate(context *gin.Context) {
	token := context.Request.Header.Get("token")
	// must be some auth verify
	switch token {
	case adminToken:
		context.Set("isAdmin", true)
	case userToken:
		context.Set("isAdmin", false)
	default:
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	context.Next()
}
```

**Варианты решения:** 

1) Использовать JWT для передачи информации о пользователе (id, является ли админом), и подписания их секретом, хранящимся на стороне сервиса. В middleware сервис стоит просто проверять подпись, после чего прокидывать информацию о конкретном пользователе. Пример

```go
func NewCheckAuth(log *zap.Logger) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("token")

			parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Warn("Unexpected signing method", zap.Any("alg", token.Header["alg"]))
					return nil, pErrors.ErrInvalidAuthToken
				}

				return []byte(viper.GetString(config.AuthKey)), nil
			})
			if err != nil {
				pHTTP.HandleError(w, r, pErrors.ErrInvalidAuthToken)
				return
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				pHTTP.HandleError(w, r, pErrors.ErrInvalidAuthToken)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims["user_id"])
			ctx = context.WithValue(ctx, ContextIsAdmin, claims["is_admin"])

			h(w, r.WithContext(ctx))
		}
	}
}
```

2) Хранить на стороне сервиса (в БД) сессии - user id : токен, и проверять наличие такой сессии.

```

## text/Бизнес-логика в handler.md
_Размер: 5028 байт_
_Кодировка: utf-8_

```
# Бизнес-логика в handler

**Ошибка:**
Выносить бизнес-логику на уровень `handler`, отказываться от слоя `service` или `usecase`

В контексте веб-приложений, обработчики `handlers` обычно отвечают за прием и отправку HTTP-запросов, извлечение данных из запроса, вызов соответствующих методов и возвращение ответа клиенту. Эти обработчики не должны содержать сложную бизнес-логику, так как это делает их менее читаемыми и трудно поддерживаемыми.

Вместо этого, бизнес-логика должна быть вынесена в отдельный слой, часто называемый слоем `usecase` или слоем `service`. Здесь содержатся структуры или функции, которые реализуют конкретные бизнес-операции.

Использование слоя `usecase` позволяет лучше структурировать код, делает его более читаемым и поддерживаемым. Также это упрощает тестирование бизнес-логики, так как мы можем написать юнит-тесты для отдельных `usecase'ов` без необходимости имитировать весь цикл обработки HTTP-запроса.

**Решение:**

- Добавляем новый слой `usecase|services` в наше приложение, переносим всю бизнес-логику в этот слой
- Пишем код в `DDD` формате, в таком случаи в слое `entity` у моделей прописываем логику. `handler` выступает в роли фасада

**Пример с ошибкой:**

```go
func GetUserBanner(
	log *slog.Logger,
	userBannerGetter UserBannerGetter,
	bannerCache BannerCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.user_banner.GetUserBanner"

		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			log.Error("error converting tagID", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}
		featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			log.Error("error converting featureID", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		useLastRevision := false
		useLastRevisionStr := r.URL.Query().Get("use_last_revision")
		if useLastRevisionStr == "true" {
			useLastRevision = true
		} else if useLastRevisionStr != "false" && useLastRevisionStr != "" {
			log.Error("Incorrect data")
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		isAdmin := r.Context().Value("isAdmin").(bool)

		var bannerContent json.RawMessage
		var bannerIsActive bool
		isCacheUsed := false
		if !useLastRevision { // Начинается бизнес-логика
			bannerContent, bannerIsActive, err = bannerCache.GetBanner(r.Context(), tagID, featureID)
			if err != nil {
				log.Error("Error fetching banner content from cache", sl.Err(err))
			} else {
				log.Info("Get data from cache, successful")
				isCacheUsed = true
			}
		}
		if useLastRevision || !isCacheUsed {
			bannerContent, bannerIsActive, err = userBannerGetter.GetUserBanner(r.Context(), tagID, featureID)
			if err != nil {
				if errors.Is(err, errs.ErrBannerNotFound) {
					log.Error("Banner is not found", sl.Err(err))
					render.JSON(w, r, response.NewError(http.StatusNotFound, "Banner is not found"))
					return
				}
				log.Error("Internal error", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Intrenal error"))
				return
			}
			err := bannerCache.SetBanner(r.Context(), tagID, featureID, &models.BannerForUser{bannerContent, bannerIsActive})
			if err != nil {
				log.Error("Error setting banner content in cache", sl.Err(err))
			} else {
				log.Info(
					"Data cached:",
					slog.Any("bannerContent", bannerContent),
					slog.Any("bannerIsActive", bannerIsActive),
					slog.Any("tagID", tagID),
					slog.Any("featureID", featureID))
			}
		}
		if !isAdmin && !bannerIsActive {
			log.Error("User have no access to inactive banner")
			render.JSON(w, r, response.NewError(http.StatusForbidden, errs.ErrUserDoesNotHaveAccess.Error()))
			return
		}
		log.Info("Successful respnose:", slog.Any("banner content", bannerContent))
		render.JSON(w, r, ResponseGet{
			response.NewSuccess(200),
			bannerContent,
		})
	}

}

```

```

## text/Варианты схемы БД.md
_Размер: 2158 байт_
_Кодировка: utf-8_

```
# Варианты схемы БД

Фича и тег у нас однозначно определяют баннер, но у одного баннера может быть несколько тегов. 

## Вариант 1

Храним в таблице баннеров id фичи и список id тегов. 

Для быстрого поиска в случае PostgreSQL можно будет использовать [GIN индекс](https://postgrespro.com/blog/pgsql/4261647) для колонки списка id тегов. 

```sql
CREATE TABLE IF NOT EXISTS banners (
     id SERIAL PRIMARY KEY,
     tag_ids integer[],
     feature_id integer,
     content jsonb,
     is_active boolean,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tag_ids ON banners USING GIN (tag_ids);
CREATE INDEX IF NOT EXISTS idx_feature_id ON banners (feature_id);
```

Проблемы:

- Гарантия и проверка уникальности пары id фичи + id тега среди баннеров — нельзя просто навесить уникальный индекс, необходимо что-то придумывать. Конечно, можно проверять наличие в транзакции на уровне приложения, однако это означает, что потенциально в БД могут существовать несогласованные данные.

## Вариант 2

Храним в отдельной таблице id тега + id фичи + id баннера, при получении делаем join. 

```sql
CREATE TABLE feature_tag_banner
(
    tag_id bigint not null,
    feature_id bigint not null, 
    banner_id bigint not null  
    primary key (tag_id, feature_id)
);

CREATE TABLE IF NOT EXISTS banners (
     id SERIAL PRIMARY KEY,
     content jsonb,
     is_active boolean,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

```

## text/Дополнительное задание версионность баннеров.md
_Размер: 2518 байт_
_Кодировка: utf-8_

```
# Дополнительное задание: версионность баннеров

## Задание

Иногда получается так, что необходимо вернуться к одной из трех предыдущих версий баннера в связи с найденной ошибкой в логике, тексте и т.д. Измените API таким образом, чтобы можно было просмотреть существующие версии баннера и выбрать подходящую версию.

## Пример решения

Создаем отдельную таблицу с версиями баннеров. При создании баннера, создаем запись еще и в таблице с версиями (в транзакции, либо асинхронно через очередь). 

Почему отдельная таблица, а не просто флаг активной версии? 

Основная нагрузка у нас будет на получение баннеров актуальной версии. Если мы предполагаем частое редактирование баннеров, то размер таблицы с 3-мя последними версиями каждого из баннеров будет гораздо больше, что может замедлить скорость выполнения запроса на получение актуального баннера/баннеров. 

### Избавление от старых версий

Нет необходимости хранить всю историю — избавляемся от старых записей.

**Сохраняем новую версию синхронно (в транзакции)**

Можно сделать крон, который будет обрабатывать таблицу с историей и удалять лишние версии (при необходимости сохранения только последних N). 

**Сохраняем новую версию асинхронно (через очередь)**

Можно удалить лишние старые версии прямо при записи новой — в транзакции проверяем количество версий в истории и удаляем самую старую для этого баннера, если требуется.

```

## text/Дополнительное задание удаление большого количества баннеров.md
_Размер: 3364 байт_
_Кодировка: utf-8_

```
# Дополнительное задание: удаление большого количества баннеров

## Задание

Добавить метод удаления баннеров по фиче или тегу, время ответа которого не должно превышать 100 мс, независимо от количества баннеров. В связи с небольшим временем ответа метода, рекомендуется ознакомиться с механизмом выполнения отложенных действий. 

## Пример решения

Что нужно учесть при решении:

1) время ответа метода удаления баннеров сильно ограничено, а операция может быть очень тяжелой - необходимо делать это асинхронно

2) баннеров, подходящих под условие, может быть так много, что их нельзя удалять одним запросом - нужно делать это батчами (небольшими пачками) 

3) если во время асинхронной обработки удаления произошел сбой или перезапуск приложения, мы все равно хотим иметь возможность обработать эту операцию - поэтому просто запустить удаление в отдельной горутине нам не подходит 

4) мы хотим удалить именно **те** баннеры, которые подходили при критерии при выполнении запроса - если пользователь инициировал удаление баннеров по тегу, а затем добавил новый баннер с таким тегом - он должен остаться

Отложенно задачу удаления можно выполнить двумя способами: 

- через отправку события на удаление в очередь и обработку очереди консюмером
- через отдельную таблицу задач, в которую будет сохраняться задача в методе удаления баннеров, а обрабатывать сами задачи будет отдельный процесс, который постоянно будет читать эту таблицу и обрабатывать задачи из нее

Примеры очередей, которые можно использовать: RabbitMQ, Kafka.

Что должно быть в сообщении / содержимом задачи? 

Удаление у нас происходит по фиче / тегу, но если мы будем использовать в сообщении именно их, есть вероятность возникновения проблемы, описанной в пункте 4. 

Поэтому, в методе удаления можно получать список id баннеров, которые мы хотим удалить, и отправлять в сообщении именно их.

```

## text/Пакет utils (1).md
_Размер: 983 байт_
_Кодировка: utf-8_

```
# Пакет utils

**Ошибка:**
Создания пакета `utils` для выноса дополнительного функционала

Пакет должен нести функционально название. Эти пакеты содержат множество несвязанных функций, поэтому их полезность трудно описать в терминах того, что предоставляет пакет.

**Решение:**
Распределить код в нужных пакетах. Данный пример можно вынести в пакет `entity`

**Пример с ошибкой:**

```go
package utils

func InitNilFieldsOfBanner(banner1 *entity.Banner, banner2 *entity.Banner) {
	if banner1.FeatureID == 0 {
		banner1.FeatureID = banner2.FeatureID
	}
	...
}

```

**Хороший пример:**

```go
package entity

func InitNilFieldsOfBanner(banner1 *Banner, banner2 *Banner) {
....
}

```

```

## text/Пакет utils.md
_Размер: 983 байт_
_Кодировка: utf-8_

```
# Пакет utils

**Ошибка:**
Создания пакета `utils` для выноса дополнительного функционала

Пакет должен нести функционально название. Эти пакеты содержат множество несвязанных функций, поэтому их полезность трудно описать в терминах того, что предоставляет пакет.

**Решение:**
Распределить код в нужных пакетах. Данный пример можно вынести в пакет `entity`

**Пример с ошибкой:**

```go
package utils

func InitNilFieldsOfBanner(banner1 *entity.Banner, banner2 *entity.Banner) {
	if banner1.FeatureID == 0 {
		banner1.FeatureID = banner2.FeatureID
	}
	...
}

```

**Хороший пример:**

```go
package entity

func InitNilFieldsOfBanner(banner1 *Banner, banner2 *Banner) {
....
}

```

```

## text/Паттерн graceful-shutdown.md
_Размер: 1882 байт_
_Кодировка: utf-8_

```
# Паттерн graceful-shutdown

**Ошибка:**
Без graceful shutdown приложение может просто отключиться, оставив открытыми соединения с базой данных, внешними сервисами или клиентами. Это может привести к утечкам ресурсов и проблемам с производительностью

**Решение:**
Реализовать паттерн в своем коде через каналы и сигналы

**Пример с ошибкой:**

```go
	// Routes
	r := gin.New()
	v1.NewRouter(r, middlewares, authService, bannerService)

	r.Run(fmt.Sprintf(":%d", config.HTTP.Port))

```

**Хороший пример:**

```go
func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

```

```

## text/Повсеместное использование транзакций.md
_Размер: 1428 байт_
_Кодировка: utf-8_

```
# Повсеместное использование транзакций

**Ошибка:**
Использование транзакции может показаться излишним —  не происходит выполнение нескольких операций, которые нужно либо выполнять вместе одновременно, либо полностью отменить

**Решение:**
Отказаться от инициализации транзакции, обходимся обычным коннектом к bd

**Пример с ошибкой:**

```go
func (s *Storage) DeleteBanners(ctx context.Context, featureID, tagID *int) (int, error) {
	const op = "repo.postgres.DeleteBanners"

	tx, err := s.db.Begin(ctx) // НЕ НУЖНА ТУТ ТРАНЗАКЦИЯ
	if err != nil {
		return 0, fmt.Errorf("%s: begin transaction %w", op, err)
	}
	defer tx.Rollback(ctx)

	query := `
        UPDATE banners
        SET deleted = true
        WHERE not deleted
        AND (feature_id = $1 OR $1 IS NULL)
        AND ($2 = ANY(tag_ids) OR $2 IS NULL);
    `

	res, err := s.db.Exec(ctx, query, featureID, tagID)
	if err != nil {
		return 0, fmt.Errorf("%s: execute context %w", op, err)
	}
	rowsAffected := res.RowsAffected()

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: commit transaction %w", op, err)
	}

	return int(rowsAffected), nil
}
```

```

## text/Покрытие тестами менее 70%.md
_Размер: 483 байт_
_Кодировка: utf-8_

```
# Покрытие тестами < 70%

**Ошибка:**

1. Не покрываем слой `repository` интеграционными тестами
2. Не покрываем слой `handler` интеграционными тестами
3. Не покрываем слой `use-case` юнит-тестами

**Решение:**

1. `repository` интеграционные тесты
2. `handler` интеграционные тесты
3. `use-case` юнит-тесты

```

## text/Построение динамических запросов к БД.md
_Размер: 4221 байт_
_Кодировка: utf-8_

```
# Построение динамических запросов к БД

**Ошибка:** 

Есть набор параметров для поиска сущностей в БД, и при этом все или часть из них опциональны. Нужно уметь строить запрос для разного подмножества этих параметров. В одном из решений был написан SQL запрос и метод репозитория для каждого из возможных таких наборов - это очень большое дублирование кода, такое тяжело поддерживать, к тому же, добавление любого нового параметра привет к добавлению еще множества таких запросов. 

Пример (не будем приводить все ветвеления, их в таком случае будет очень много): 

```go
switch {
	// nothing is provided, return empty slice
	case params.FeatureId == nil && params.TagId == nil:
		break

	// only feature id is provided
	case params.FeatureId != nil && params.TagId == nil &&
		params.Limit == nil && params.Offset == nil:
		res, err := s.repo.GetBannersByFeature(r.Context(), params)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			ErrorHandlerFunc(w, r, err)
			return
		}
		response = res

	// feature id + limit
	case params.FeatureId != nil && params.Limit != nil &&
		params.TagId == nil && params.Offset == nil:
		res, err := s.repo.GetBannersByFeatureWithLimit(r.Context(), params)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			ErrorHandlerFunc(w, r, err)
			return
		}
		response = res

	// feature id + offset
	case params.FeatureId != nil && params.Offset != nil &&
		params.TagId == nil && params.Limit == nil:
		res, err := s.repo.GetBannersByFeatureWithOffset(r.Context(), params)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			ErrorHandlerFunc(w, r, err)
			return
		}
		response = res
		
	// ... и так далее
	
	}
```

**Варианты решения:** 

1) Можно составить динамический запрос вручную, формируя строку запроса и список параметров. Пример:

```go
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 4)
	if params.FeatureID > 0 {
		conditions = append(conditions, fmt.Sprintf("feature_id = $%d", len(args)+1))
		args = append(args, params.FeatureID)
	}
	if params.TagID > 0 {
		conditions = append(conditions, fmt.Sprintf("tag_id = $%d", len(args)+1))
		args = append(args, params.TagID)
	}

	var conditionPart string
	if len(conditions) > 0 {
		condition := strings.Join(conditions, " AND ")
		conditionPart = fmt.Sprintf(`
WHERE b.id IN (SELECT banner_id
               FROM banner_references
               WHERE %s)`, condition)
	}

	var limitPart string
	if params.Limit > 0 {
		limitPart = fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, params.Limit)
	}
	if params.Offset > 0 {
		limitPart += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, params.Offset)
	}

	cmd := fmt.Sprintf(listCmd, conditionPart, limitPart)

	rows, err := r.pool.Query(ctx, cmd, args...)
```

2) Можно использовать библиотеки-билдеры запросов, например [github.com/Masterminds/squirrel](http://github.com/Masterminds/squirrel). 

Пример использования (здесь забыли про limit и offset, но как пример использования библиотеки вполне подходит)

```go
psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	searchBannersID := psql.Select("ID").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id")

	if bannerParams.FeatureId != 0 {
		searchBannersID = searchBannersID.Where("feature_id = ?", bannerParams.FeatureId)
	}
	if bannerParams.TagId != 0 {
		searchBannersID = searchBannersID.Where("tag_id = ?", bannerParams.TagId)
	}

	bannersQuery, args, err = searchBannersID.ToSql()

```

```

## text/Проверка на nil.md
_Размер: 586 байт_
_Кодировка: utf-8_

```
# Проверка на nil

**Ошибка:**
Не проверять указатели на nil (приводит к панике)

**Решение:**
Добавлять проверку на nil

**Пример с ошибкой:**

```go
func InitNilFieldsOfBanner(banner1 *entity.Banner, banner2 *entity.Banner) {
	if banner1.FeatureID == 0 {
		banner1.FeatureID = banner2.FeatureID
	}
	...
}

```

**Хороший пример:**

```go
func InitNilFieldsOfBanner(banner1 *entity.Banner, banner2 *entity.Banner) {
	if banner1 == nil && banner2 == nil {
		return
	}
	...
}

```

```

## text/Функции стандартной библиотеки.md
_Размер: 624 байт_
_Кодировка: utf-8_

```
# Функции стандартной библиотеки

**Ошибка:**
Писать обертки для слайсов или строк, вместо использования функции стандартной библиотеки

**Решение:**
Использовать пакет `slices` метод `Equals`

**Пример с ошибкой:**

```go
func Equals[T comparable](s1 []T, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}
```

**Хороший пример:**

```go
slices.Equals(s1, s2)
```

```
