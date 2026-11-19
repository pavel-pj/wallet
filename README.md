[![tests](https://github.com/pavel-pj/wallet/actions/workflows/tests.yml/badge.svg)](https://github.com/pavel-pj/wallet/actions/workflows/tests.yml)
[![lint](https://github.com/pavel-pj/wallet/actions/workflows/linter.yaml/badge.svg)](https://github.com/pavel-pj/wallet/actions/workflows/linter.yaml)
[![integration-tests](https://github.com/pavel-pj/wallet/actions/workflows/integration-tests.yaml/badge.svg)](https://github.com/pavel-pj/wallet/actions/workflows/integration-tests.yaml)
# Wallet Service

Сервис кошельков: gRPC, HTTP/JSON API и Swagger UI, база — PostgreSQL.
 
 

## Запуск

```bash
cp config.env.example config.env
```

Заполнить `config.env` и запустить:

```bash
make dev-build   # первый запуск
make dev-up      # обычный запуск
make dev-logs    # логи приложения
make dev-down    # остановить (удаляет данные БД!)
```

Полезные команды:

```bash
make tests-cover               # юнит-тесты, отчет в coverage.html
make tests-integration-cover   # интеграционные тесты, отчет в coverage.integration.html
make lint                      # линтер
```
 
При старте автоматически применяются миграции, затем поднимается приложение.

После запуска:

- HTTP API: `http://localhost:50470`
- gRPC: `localhost:50460`
- Swagger: `http://localhost:50480`


## Порты

| Порт | Что |
|---|---|
| `50460` | gRPC |
| `50470` | HTTP/JSON API |
| `50480` | Swagger UI |
| `6001` | PostgreSQL снаружи (пример; задаётся `PG_PORT`) |

## Swagger

Открыть `http://localhost:50480` — там описаны HTTP-методы и можно отправлять запросы.

 ## API

| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/create` | создать кошелёк |
| POST | `/api/v1/wallet` | пополнить или списать |
| GET | `/api/v1/{uuid}` | получить баланс кошелька |

Создать кошелёк:

```bash
curl -X POST http://localhost:50470/api/v1/create \
  -H 'Content-Type: application/json' \
  -d '{"amount": 2000}'
```

Пополнить:

```bash
curl -X POST http://localhost:50470/api/v1/wallet \
  -H 'Content-Type: application/json' \
  -d '{"wallet_id":"<uuid>","operation_type":"DEPOSIT","amount":544}'
```

Списать:

```bash
curl -X POST http://localhost:50470/api/v1/wallet \
  -H 'Content-Type: application/json' \
  -d '{"wallet_id":"<uuid>","operation_type":"WITHDRAW","amount":544}'
```

`operation_type` — `DEPOSIT` или `WITHDRAW`. `wallet_id` — UUID кошелька, полученный при создании.

Получить баланс:

```bash
curl http://localhost:50470/api/v1/<uuid>
```

Пример ответа:


```bash
{
  "wallet": {
    "uuid": "7d8bf21f-ada5-4b3b-9e6a-fd4381f7348e",
    "amount": "2000",
    "createdAt": "2026-11-19T12:00:00Z",
    "updatedAt": "2026-11-19T12:00:00Z"
  }
}
```
 
 
## Конфигурация

Все настройки — в `config.env`:

| Переменная | Значение |
|---|---|
| `POSTGRES_DB` | имя БД |
| `POSTGRES_USER` | пользователь БД |
| `POSTGRES_PASSWORD` | пароль БД |
| `DB_HOST` | хост БД для приложения (`localhost`) |
| `DB_PORT` | порт БД для приложения (`6001`) |
| `PG_PORT` | внешний порт PostgreSQL (`6001`) |
| `SERVER_PORT` | порт gRPC (`50460`) |
| `SERVER_HTTP_PORT` | порт HTTP API (`50470`) |
| `SERVER_SWAGGER_PORT` | порт Swagger (`50480`) |

Пример:

```dotenv
POSTGRES_DB=db_wallet
POSTGRES_USER=wallet
POSTGRES_PASSWORD=change_me
DB_HOST=localhost
DB_PORT=6001
PG_PORT=6001

SERVER_PORT=50460
SERVER_HTTP_PORT=50470
SERVER_SWAGGER_PORT=50480
SERVER_SWAGGER_HOST=localhost
```

## Параметры нагрузочного теста

Для тестирование используется команда(запускать из другого терминала при запущенном сервере) :

| Параметр | Описание |
|---|---|
| `wallet` | UUID кошелька (сначала создать в БД через `Store`) |
| `rate` | RPS — целевое количество запросов в секунду |
| `duration` | Продолжительность теста, например `10s` |
| `workers` | Количество воркеров для тестирования |
| `op` | Операция: `DEPOSIT` или `WITHDRAW` |
| `amount` | Сумма пополнения или списания |

Пример запуска:

```bash
go run ./cmd/loadtest \
  -wallet <uuid> \
  -rate 1000 \
  -duration 10s \
  -workers 200 \
  -op DEPOSIT \
  -amount 1
```

Примерный результат тестирования:

```bash
ok=4788 failed=0 5xx=0 elapsed=10.651913662s actual_rps=449
```
