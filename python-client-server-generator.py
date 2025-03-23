import sys

if len(sys.argv) != 3:
    raise Exception("Se deben recibir exactamente 2 parametros.")

archivo_salida = sys.argv[1]

try:
    cantidad_clientes = int(sys.argv[2])
except:
    raise Exception("El segundo parametro debe ser un entero.")

if cantidad_clientes < 0:
    raise Exception("Se debe recibir como minimo 1 cliente.")

string_compose = """name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    volumes:
      - ./server/config.ini:/config.ini
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
      - NOMBRE="Santiago"
      - APELLIDO="Lorca"
      - DOCUMENTO=30904465
      - NACIMIENTO=1999-03-17
      - NUMERO=7574
    volumes:
      - ./client/config.yaml:/config.yaml
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
try:
    archivo = open(archivo_salida, "w")
except:
    raise Exception("Error al intentar abrir el archivo en escritura.")
archivo.write(string_compose)
archivo.close()

print("Archivo docker-compose.yaml generado con éxito")
