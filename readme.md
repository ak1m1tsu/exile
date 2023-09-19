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
KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://broker:29092,PLAINTEXT_HOST://localhost:9092
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
```

Доалее используя команду `make up` запустить докер контейнеры.

## `Makefile` команды

| Команда 	| Описание 	|
|---	|---	|
| `make up` 	| Запускает скрипт `dockerup.sh`, который билдит все докер контейнеры. 	|
| `make down` 	| Останавливает запущенные контейнеры 	|
| `make gen` 	| Генерирует моки для интерфейсов, используя [mockery](https://github.com/vektra/mockery) 	|
| `make kafkaseed` 	| Создает топики **FIO** и **FIO_FAILED**  	|
| `make tests` 	| Запускает unit-тесты 	|
