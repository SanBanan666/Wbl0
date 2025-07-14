# 🚀 Order Service - Микросервис обработки заказов

Демонстрационный микросервис на Go для обработки заказов из Kafka, сохранения в PostgreSQL, кеширования в памяти и предоставления HTTP API с веб-интерфейсом.

## 📋 Содержание

- [Архитектура](#архитектура)
- [Технологии](#технологии)
- [Структура проекта](#структура-проекта)
- [Установка и запуск](#установка-и-запуск)
- [Использование](#использование)
- [API документация](#api-документация)
- [Разработка](#разработка)
- [Мониторинг](#мониторинг)

## 🏗️ Архитектура

Проект реализует микросервисную архитектуру с следующими компонентами:

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Frontend  │    │   Kafka UI  │    │   Script    │
│  (Nginx)    │    │  (Port 8080)│    │  Producer   │
│ (Port 8082) │    └─────────────┘    └─────────────┘
└─────────────┘           │                   │
         │                │                   │
         ▼                ▼                   ▼
┌─────────────────────────────────────────────────────────┐
│                    Order Service                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   HTTP API  │  │   Kafka     │  │   Cache     │     │
│  │ (Port 8081) │  │ Consumer    │  │ (In-Memory) │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
         │                │                   │
         ▼                ▼                   ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ PostgreSQL  │    │   Kafka     │    │  Zookeeper  │
│ (Port 5432) │    │ (Port 9092) │    │ (Port 2181) │
└─────────────┘    └─────────────┘    └─────────────┘
```

### Компоненты системы:

1. **Order Service (Go)** - основной микросервис
   - HTTP API для получения заказов
   - Kafka Consumer для обработки сообщений
   - In-Memory кеш для ускорения запросов
   - PostgreSQL репозиторий для хранения данных

2. **PostgreSQL** - основная база данных
   - Хранение заказов, доставки, платежей и товаров
   - Автоматические миграции при запуске

3. **Kafka + Zookeeper** - брокер сообщений
   - Асинхронная обработка заказов
   - Масштабируемость и отказоустойчивость

4. **Nginx** - веб-сервер для фронтенда
   - Статическая раздача HTML/CSS/JS
   - Проксирование API запросов

5. **Kafka UI** - веб-интерфейс для управления Kafka
   - Просмотр топиков и сообщений
   - Мониторинг состояния брокера

## 🛠️ Технологии

### Backend
- **Go 1.23** - основной язык разработки
- **PostgreSQL 15** - реляционная база данных
- **Apache Kafka 7.4** - брокер сообщений
- **Sarama** - Go клиент для Kafka
- **lib/pq** - драйвер PostgreSQL

### Frontend
- **HTML5/CSS3** - современный веб-интерфейс
- **Vanilla JavaScript** - интерактивность
- **Nginx** - веб-сервер

### Infrastructure
- **Docker & Docker Compose** - контейнеризация
- **Alpine Linux** - легковесные образы

## 📁 Структура проекта

```
wildL0/
├── Wbl0/                          # Основной код Go
│   ├── cmd/
│   │   ├── api/                   # HTTP API сервер
│   │   │   ├── main.go           # Точка входа
│   │   │   └── Dockerfile        # Образ для API
│   │   └── migr/                 # Миграции БД
│   │       ├── main.go
│   │       └── Dockerfile
│   ├── internal/
│   │   ├── application/          # Слой приложения
│   │   │   ├── dto/             # Data Transfer Objects
│   │   │   ├── handlers/        # Обработчики сообщений
│   │   │   ├── interfaces/      # Интерфейсы
│   │   │   └── services/        # Бизнес-логика
│   │   ├── domain/              # Доменная модель
│   │   │   └── entities/        # Сущности
│   │   ├── infrastructure/      # Инфраструктурный слой
│   │   │   ├── consumers/       # Kafka потребители
│   │   │   ├── migrations/      # Миграции БД
│   │   │   └── repositories/    # Репозитории
│   │   └── presentation/        # Слой представления
│   │       └── controllers/     # HTTP контроллеры
│   ├── pkg/                     # Общие пакеты
│   │   ├── dbconnections/       # Подключения к БД
│   │   └── logger/              # Логирование
│   └── db/
│       └── migrations/          # SQL миграции
├── frontend/                     # Веб-интерфейс
│   └── index.html               # Главная страница
├── scripts/                      # Скрипты
│   └── send_test_order.go       # Тестовый producer
├── docker-compose.yml           # Конфигурация Docker
├── nginx.conf                   # Конфигурация Nginx
├── go.mod                       # Go модули
├── go.sum                       # Go зависимости
└── README.md                    # Документация
```

## 🚀 Установка и запуск

### Предварительные требования

- Docker Desktop
- Go 1.23+ (для локальной разработки)
- Git

### Быстрый запуск

1**Запуск всех сервисов:**
```bash
docker-compose up -d
```

2**Проверка статуса:**
```bash
docker-compose ps
```

3**Отправка тестовых данных:**
```bash
go run scripts/send_test_order.go
```

### Доступные сервисы

После запуска будут доступны следующие сервисы:

| Сервис | URL | Описание |
|--------|-----|----------|
| Frontend | http://localhost:8082 | Веб-интерфейс поиска заказов |
| API | http://localhost:8081 | HTTP API для заказов |
| Kafka UI | http://localhost:8080 | Управление Kafka |
| PostgreSQL | localhost:5432 | База данных |
| Kafka | localhost:9092 | Брокер сообщений |

## 📖 Использование

### Веб-интерфейс

1. Откройте http://localhost:8082
2. Введите ID заказа в поле поиска
3. Нажмите "Найти заказ"
4. Просмотрите детальную информацию о заказе

### HTTP API

#### Получение заказа по ID
```bash
GET http://localhost:8081/order/{order_uid}
```

**Пример запроса:**
```bash
curl http://localhost:8081/order/b563feb7b2b84b6test
```

**Пример ответа:**
```json
{
  "order": {
    "order_uid": "b563feb7b2b84b6test",
    "track_number": "WBILMTESTTRACK",
    "entry": "WBIL",
    "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
    },
    "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
    },
    "items": [
      {
        "chrt_id": 9934930,
        "track_number": "WBILMTESTTRACK",
        "price": 453,
        "rid": "ab4219087a764ae0btest",
        "name": "Mascaras",
        "sale": 30,
        "size": "0",
        "total_price": 317,
        "nm_id": 2389212,
        "brand": "Vivienne Sabo",
        "status": 202
      }
    ],
    "locale": "en",
    "internal_signature": "",
    "customer_id": "test",
    "delivery_service": "meest",
    "shardkey": "9",
    "sm_id": 99,
    "date_created": "2021-11-26T06:22:19Z",
    "oof_shard": "1"
  }
}
```

#### Проверка здоровья сервиса
```bash
GET http://localhost:8081/health
```

### Отправка тестовых заказов

```bash
go run scripts/send_test_order.go
```

Этот скрипт отправляет 3 тестовых заказа в Kafka:
- `b563feb7b2b84b6test`
- `test-order-1`
- `test-order-2`

## 🔧 API документация

### Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/order/{order_uid}` | Получить заказ по ID |
| GET | `/health` | Проверка здоровья сервиса |

### Коды ответов

| Код | Описание |
|-----|----------|
| 200 | Успешный запрос |
| 404 | Заказ не найден |
| 500 | Внутренняя ошибка сервера |

### CORS

API поддерживает CORS для веб-приложений:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## 💻 Разработка

### Локальная разработка

1. **Запуск зависимостей:**
```bash
docker-compose up -d postgres kafka zookeeper
```

2. **Настройка переменных окружения:**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=orders_db
export DB_SSLMODE=disable
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=orders
export KAFKA_GROUP_ID=order-service-group
export HTTP_PORT=8081
```

3. **Запуск сервиса:**
```bash
go run Wbl0/cmd/api/main.go
```

### Структура кода

Проект следует принципам Clean Architecture:

- **Domain Layer** - бизнес-сущности и правила
- **Application Layer** - бизнес-логика и сервисы
- **Infrastructure Layer** - внешние зависимости (БД, Kafka)
- **Presentation Layer** - HTTP контроллеры

### Добавление новых функций

1. **Новый endpoint:**
   - Добавить метод в контроллер
   - Зарегистрировать маршрут в main.go

2. **Новая бизнес-логика:**
   - Создать сервис в application/services
   - Определить интерфейс в application/interfaces

3. **Новые данные:**
   - Добавить миграцию в db/migrations
   - Обновить сущности в domain/entities

## 📊 Мониторинг

### Логи

Просмотр логов сервисов:

```bash
# Логи Order Service
docker logs order-service-app

# Логи PostgreSQL
docker logs order-service-postgres

# Логи Kafka
docker logs order-service-kafka
```

### Kafka UI

Откройте http://localhost:8080 для:
- Просмотра топиков
- Мониторинга сообщений
- Проверки состояния consumer групп

### Метрики производительности

- **Кеш-хиты:** Логи показывают `Order found in cache`
- **Время ответа API:** Измеряется через curl
- **Обработка Kafka:** Логи `Received message` и `processed successfully`

## 🔍 Устранение неполадок

### Частые проблемы

1. **Порт 80 занят Apache2:**
   - Остановите Apache2: `sudo systemctl stop apache2`
   - Или измените порт Nginx в docker-compose.yml

2. **Ошибка подключения к БД:**
   - Проверьте переменные окружения
   - Убедитесь, что PostgreSQL запущен

3. **Kafka не подключается:**
   - Проверьте, что Zookeeper запущен
   - Убедитесь в правильности адресов брокеров

4. **CORS ошибки:**
   - Проверьте, что CORS middleware добавлен
   - Убедитесь в правильности заголовков

### Команды диагностики

```bash
# Проверка статуса контейнеров
docker-compose ps

# Проверка логов
docker-compose logs

# Перезапуск сервиса
docker-compose restart order-service

# Полная пересборка
docker-compose down -v
docker-compose up -d --build
```

## 📝 Лицензия

Этот проект создан в демонстрационных целях.

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request# 
Wbl0