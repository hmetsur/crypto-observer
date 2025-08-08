# Crypto Observer

Crypto Observer — это сервис для отслеживания цен криптовалют с периодическим обновлением и сохранением в базу данных.

## Возможности
- Добавление криптовалюты в список отслеживания
- Удаление криптовалюты из списка
- Получение цены криптовалюты на определённый момент времени
- Интеграция с API CoinGecko
- Хранение данных в PostgreSQL
- Логирование с уровнями info/error

## API эндпоинты

### Добавить валюту
POST /currency/add
Content-Type: application/json

{
  "symbol": "btc",
  "period": 5
}

### Удалить валюту
POST /currency/remove
Content-Type: application/json

{
  "symbol": "btc"
}

### Получить цену
GET /price?symbol=btc&timestamp=1691500000

## Запуск

### Локально
go run ./cmd

### Через Docker
docker-compose up --build

## Переменные окружения
- DB_DSN — строка подключения к PostgreSQL (например, postgres://user:pass@localhost:5432/crypto?sslmode=disable)
- LOG_LEVEL — уровень логирования (info, error, debug)

## Тестирование
Для запуска всех тестов с отчётом покрытия:
go test ./... -cover

## Примеры запросов

### Добавить валюту в отслеживание
curl -X POST "http://localhost:8080/currency/add" \
     -H "Content-Type: application/json" \
     -d '{"symbol":"btc","period":5}'

### Удалить валюту из мониторинга
curl -X POST "http://localhost:8080/currency/remove" \
     -H "Content-Type: application/json" \
     -d '{"symbol":"btc"}'

### Получить последнюю цену
curl "http://localhost:8080/price?symbol=btc"

### Пример ответа
{
    "symbol": "btc",
    "timestamp": 1723112000,
    "price": 29150.32
}

