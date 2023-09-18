docker compose exec kafka \
  kafka-topics --create \
    --topic FIO \
    --bootstrap-server localhost:9092 \
    --replication-factor 1 \
    --partitions 1

docker compose exec kafka \
  kafka-topics --create \
    --topic FIO_FAILED \
    --bootstrap-server localhost:9092 \
    --replication-factor 1 \
    --partitions 1
