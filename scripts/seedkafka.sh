docker compose exec broker \
  kafka-topics --create \
    --topic FIO \
    --bootstrap-server broker:9092 \
    --replication-factor 1 \
    --partitions 1

docker compose exec broker \
  kafka-topics --create \
    --topic FIO_FAILED \
    --bootstrap-server broker:9092 \
    --replication-factor 1 \
    --partitions 1
