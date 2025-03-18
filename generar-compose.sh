#!/bin/bash

# Revisamos que hayan llegado dos parametros, en caso de que no, frenamos
# antes de levantar python.
if [ "$#" -ne 2 ]; then
    echo "Error: Se requieren exactamente 2 parámetros."
    exit 1
fi

echo "Archivo de salida: $1"
echo "Cantidad de clientes: $2"

python3 python-client-server-generator.py "$1" "$2"
