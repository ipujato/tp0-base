import sys

# Revisamos cuantos argumentos llegan 
if len(sys.argv) != 3:
    raise Exception("Se deben recibir exactamente 2 parametros.")

archivo_salida = sys.argv[1]

# Comprobamos que clientes sea un entero
try:
    cantidad_clientes = int(sys.argv[2])
except:
    raise Exception("El segundo parametro debe ser un entero.")

string_compose = """
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
"""

for i in range(1, cantidad_clientes + 1):
    string_compose += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={i}
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
"""

string_compose += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

archivo = open(archivo_salida, "w")
archivo.write(string_compose)
archivo.close()

print("Archivo docker-compose.yaml generado con éxito")
