#!/bin/bash

#https://hub.docker.com/r/subfuzion/netcat

echo "Mensaje ej3"
RESPUESTA=$(docker run --rm --platform linux/amd64 --network=testing_net -i subfuzion/netcat server 12345)

if [ "$RESPUESTA" = "Mensaje ej3" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi