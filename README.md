# Wallet Service

Сервис кошельков: gRPC, HTTP/JSON API и Swagger UI, база — PostgreSQL.
 
## Порты

| Порт | Что |
|---|---|
| `50460` | gRPC |
| `50470` | HTTP/JSON API |
| `50480` | Swagger UI |
| `6001` | PostgreSQL снаружи (пример; задаётся `PG_PORT`) |

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

При старте автоматически применяются миграции, затем поднимается приложение.

После запуска:

- HTTP API: `http://localhost:50470`
- gRPC: `localhost:50460`
- Swagger: `http://localhost:50480`

## API

| Метод | Путь | Назначение |
|---|---|---|
| POST | `/api/v1/create` | создать кошелёк |
| POST | `/api/v1/wallet` | пополнить или списать |

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

gRPC-сервис: `wallet_v1.WalletV1` (`Store`, `Wallet`), порт `50460`.

## Swagger

Открыть `http://localhost:50480` — там описаны HTTP-методы и можно отправлять запросы.
