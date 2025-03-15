#!/bin/bash

# Revisamos que hayan llegado dos parametros, en caso de que no, frenamos
# antes de levantar python.
if [ "$#" -ne 2 ]; then
    echo "Error: Se requieren exactamente 2 parámetros."
    exit 1
fi

# Impresion para controlar los parametros
echo "Archivo de salida: $1"
echo "Cantidad de clientes: $2"

python3 python-client-server-generator.py "$1" "$2"

# 
if [ "$?" -ne 0 ]; then
    echo "Error: Falló la ejecución de python-client-server-generator.py"
    exit 1
fi