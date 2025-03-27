#!/bin/bash

TEST_MESSAGE="Mensaje comprobacion ej3"

RESPUESTA=$(echo "$TEST_MESSAGE" | docker run --rm --platform linux/amd64 --network=tp0_testing_net alpine sh -c "nc server 12345 -w 3")

if [ "$RESPUESTA" = "$TEST_MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi