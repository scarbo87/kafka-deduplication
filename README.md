```
docker run --rm -it --net kafka-deduplication_default confluentinc/cp-kafkacat kafkacat -b broker:29092 -L

docker run -ti --net kafka-deduplication_default confluentinc/cp-kafkacat kafkacat -b broker:29092 -C -K: -f '\nKey (%K bytes): %k\t\nValue (%S bytes): %s\n\Partition: %p\tOffset: %o\n--\n' -t test

docker run -ti --net kafka-deduplication_default confluentinc/cp-kafkacat kafkacat -b broker:29092 -C -K: -f '\nKey (%K bytes): %k\t\nValue (%S bytes): %s\n\Partition: %p\tOffset: %o\n--\n' -t __consumer_offsets

docker run --rm -i --net kafka-deduplication_default confluentinc/cp-kafkacat kafkacat -b broker:29092 -t test -K: \
-P <<EOF
9:9
EOF

go mod tidy -v
go mod vendor
go build -mod=vendor -v -o main ./*.go
```