#!/bin/bash

RESPUESTA=$(echo "Mensaje ej3" | docker run --rm --platform linux/amd64 --network=tp0_testing_net -i subfuzion/netcat -w 3 server 12345)

if [ "$RESPUESTA" = "Mensaje ej3" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi