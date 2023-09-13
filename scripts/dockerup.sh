#!/bin/bash
source .env

docker compose up -d

echo "Waiting for setup docker images... (15s)"
sleep 15s

docker compose exec broker \
  kafka-topics --create \
    --topic fio \
    --bootstrap-server localhost:9092 \
    --replication-factor 1 \
    --partitions 1

docker compose exec broker \
  kafka-topics --create \
    --topic fio_failed \
    --bootstrap-server localhost:9092 \
    --replication-factor 1 \
    --partitions 1
