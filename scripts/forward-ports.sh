#!/bin/bash

echo "forwarding ports needed for running integration tests"


fuser -k 9092/tcp || true
fuser -k 4566/tcp || true

kubectl port-forward svc/localstack 4566:4566 &
kubectl port-forward svc/kafka 9092:9092 &


until curl -s http://localhost:4566 > /dev/null; do
    sleep 1
done