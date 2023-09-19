# Effective Mobile test task

## Суть задания

Реализовать сервис, который будет получать поток ФИО, из открытых апи обогащать
ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в
БД. По запросу выдавать инфу о найденных людях. Необходимо реализовать следующее:

1. Сервис слушает очередь кафки **FIO**, в котором приходит информация с ФИО в
формате:

```json
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich" // необязательно
}
```
2. В случае некорректного сообщения, обогатить его причиной ошибки (нет
обязательного поля, некорректный формат...) и отправить в очередь кафки
**FIO_FAILED**
3. Корректное сообщение обогатить:
    - Возрастом - https://api.agify.io/?name=Dmitriy
    - Полом - https://api.genderize.io/?name=Dmitriy
    - Национальностью - https://api.nationalize.io/?name=Dmitriy
4. Обогащенное сообщение положить в БД postgres (структура БД должна быть создана
путем миграций)
5. Выставить REST методы:
    - Для получения данных с различными фильтрами и пагинацией
    - Для добавления новых людей
    - Для удаления по идентификатору
    - Для изменения сущности
6. Выставить GraphQL методы аналогичные REST
7. Предусмотреть кэширование данных в redis
8. Покрыть код логами
9. Покрыть бизнес-логику unit-тестами
10. Вынести все конфигурационные данные в .env

## Как запустить

Создать конфигурационный файл `.env` и заполнить его следующим содержимым:

```shell
# Postgres configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgrespwd
POSTGRES_DB=emdb
# Kafka configuration
KAFKA_BROKER_ID=1
KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://broker:29092,PLAINTEXT_HOST://broker:9092
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0
KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1
KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1
KAFKA_PROCESS_ROLES=broker,controller
KAFKA_NODE_ID=1
KAFKA_CONTROLLER_QUORUM_VOTERS=1@broker:29093
KAFKA_LISTENERS=PLAINTEXT://broker:29092,CONTROLLER://broker:29093,PLAINTEXT_HOST://0.0.0.0:9092
KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT
KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER
KAFKA_LOG_DIRS=/tmp/kraft-combined-logs
CLUSTER_ID=MkU3OEVBNTcwNTJENDM2Qk
# Redis configuration
REDIS_PASSWORD=redispwd
# Service configuration
SERVICE_ENV=development
SERVICE_KAFKA_GROUP_ID="effective-mobile"
SERVICE_KAFKA_BOOTSTRAP_SERVERS="broker:9092"
SERVICE_KAFKA_AUTO_OFFSET_RESET=earliest
SERVICE_KAFKA_PRODUCER_TOPIC=FIO_FAILED
SERVICE_KAFKA_CONSUMER_TOPICS=FIO
SERVICE_KAFKA_TIMEOUT=100ms
# API configuration
API_ENV=development
API_PORT=5555
API_IDLE_TIMEOUT=30s
API_READ_TIMEOUT=5s
API_WRITE_TIMEOUT=5s
API_KAFKA_BOOTSTRAP_SERVERS="broker:9092"
API_KAFKA_PRODUCER_TOPIC=FIO
# URLs
DATABASE_URL="postgres://postgres:postgrespwd@postgres:5432/emdb?sslmode=disable"
CACHE_URL="redis://default:redispwd@redis:6379/0"
```

Далее используем данную комбинацию команд, чтобы запустить контейнеры и создать топики в кафке

```shell
make up && make seedkafka
```

## `Makefile` команды

| Команда 	| Описание 	|
|---	|---	|
| `make up` 	| Запускает скрипт `dockerup.sh`, который билдит все докер контейнеры. 	|
| `make down` 	| Останавливает запущенные контейнеры 	|
| `make gen` 	| Генерирует моки для интерфейсов, используя [mockery](https://github.com/vektra/mockery) 	|
| `make seedkafka` 	| Создает топики **FIO** и **FIO_FAILED**  	|
| `make tests` 	| Запускает unit-тесты 	|

## Endpoints

```http
GET /person - список людей
GET /person/{id} - конкретная персона
DELETE /person/{id} - удалить конкретную персону
PATCH /person/{id} - обновить конкретную персону
POST /person - добавить новую персону
POST /person/graphql - graphql запросы по персонам
```
