#!/bin/bash

if ! docker image inspect alpine >/dev/null 2>&1; then
  echo "Alpine no encontrado. Descargándolo..."
  docker pull alpine
fi

TEST_MESSAGE="Mensaje comprobacion ej3"

RESPUESTA=$(echo "$TEST_MESSAGE" | docker run --rm --platform linux/amd64 --network=tp0_testing_net alpine sh -c "nc server 12345 -w 3")

if [ "$RESPUESTA" = "$TEST_MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi