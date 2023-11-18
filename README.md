# Exile

## Technologies

- Redis (cache)
- PostgreSQL (database)
- Apache Kafka (communication between services)
- GraphQL (API gateway endpoints)
- Docker (deploy services, database, etc.)
- External APIs
    - https://api.agify.io
    - https://api.genderize.io
    - https://api.nationalize.io

## Design diagram

![](./assets/exile-system-design.svg)

## Services

### Person service

A service that receives a stream of full name, enriches the response with the most probable age, gender and nationality using open APIs and stores the data in the database.

### API Gateway

### Get list of persons

```shell
curl http://localhost:5555/person
```

**Response**

```json
{
  "status": "OK",
  "people": [
    {
      "ID": "05dd6483-1938-4d8b-9a45-7f61a69ad377",
      "Name": "Ivan",
      "Surname": "Ivanov",
      "Patronymic": "",
      "Age": 54,
      "Gender": "male",
      "Nationality": "HR",
      "IsDeleted": false
    }
  ]
}
```

### Get person

```shell
curl http://localhost:5555/person/<id>
```

**Response**

```json
{
  "status": "OK",
  "person": {
    "ID": "05dd6483-1938-4d8b-9a45-7f61a69ad377",
    "Name": "Ivan",
    "Surname": "Ivanov",
    "Patronymic": "",
    "Age": 54,
    "Gender": "male",
    "Nationality": "HR",
    "IsDeleted": false
  }
}
```

### Create new person

```shell
curl -X POST --data '{"name":"Ivan", "surname":"Ivanov", "patronymic":"Ivanovich"}' http://localhost:5555/person
```

**Response**

```json
{
    "status": "OK",
    "person": {
        "ID": "05dd6483-1938-4d8b-9a45-7f61a69ad377",
        "Name": "Ivan",
        "Surname": "Ivanov",
        "Patronymic": "Ivanovich",
        "Age": 54,
        "Gender": "male",
        "Nationality": "HR",
        "IsDeleted": false
    }
}
```

### Update a person

```shell
curl -X PATCH --data '{"name":"Roman","surname":"Kravchuk","age":21}' http://localhost:5555/person/<id>
```

**Response**

```json
{
  "status": "OK",
  "person": {
    "ID": "05dd6483-1938-4d8b-9a45-7f61a69ad377",
    "Name": "Roman",
    "Surname": "Kravchuk",
    "Patronymic": "Ivanovich",
    "Age": 21,
    "Gender": "male",
    "Nationality": "HR",
    "IsDeleted": false
  }
}
```

### Delete a person

```shell
curl -X DELETE http://localhost:5555/person/<id>
```

**Response**

```json
{
    "status": "OK"
}
```

### GraphQL endpoints

```shell
curl -X POST http://localhost:5555/person/graphql
```

## How to run?

Create config file `.env`:

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
SERVICE_KAFKA_GROUP_ID="exile"
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

Next, use this commands to start containers and create Kafka topics

```shell
make up && make seedkafka
```

## `Makefile` commands

| Команда 	| Описание 	|
|---	|---	|
| `make up` 	| Запускает скрипт `dockerup.sh`, который билдит все докер контейнеры. 	|
| `make down` 	| Останавливает запущенные контейнеры 	|
| `make gen` 	| Генерирует моки для интерфейсов, используя [mockery](https://github.com/vektra/mockery) 	|
| `make seedkafka` 	| Создает топики **FIO** и **FIO_FAILED**  	|
| `make tests` 	| Запускает unit-тесты 	|
