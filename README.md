# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación. Respetar las entradas de log planteadas en los ejercicios, pues son las que se chequean en cada uno de los tests.

La corrección personal tendrá en cuenta la calidad del código entregado y casos de error posibles, se manifiesten o no durante la ejecución del trabajo práctico. Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).

# Ejecucion y consideraciones

## Parte 1: Introducción a Docker

### Ej1
Para ejecutarlo primero debí correr `chmod +x generar-compose.sh` debido a que no tenia permisos de ejecución. Es posible que no sea necesario. Luego se corre como dice el enunciado, desde raíz: `./generar-compose.sh docker-compose-dev.yaml 5`.

Como casos de error tuve en cuenta: parámetros insuficientes, cantidades invalidas de clientes (negativos), archivos inválidos y parámetros inválidos. 
Comprobé en bash la cantidad de parámetros para no levantar python innecesariamente ya que es una comprobación fácil y rápida realizable en bash sin problemas.


### Ej2
Se corre con los make provistos por la catedra.
Para ello simplemente se modifica el generar-compose para que se incluya como volume las configuraciones del cliente y el servidor respectivamente.
Con esto se puede facilmente modificarlos con sus configuraciones que viven dentro del volume.


### Ej3 
Para ejecutarlo primero debí correr `chmod +x validar-echo-server.sh` debido a que no tenia permisos de ejecución. Es posible que sea necesario. Luego se corre como dice el enunciado, desde raíz: `./validar-echo-server.sh`.
Lo que hace es correr un pequeño script de docker que crea una red dentro de un contenedor para hacer un echo en la red sobre un netcat. Luego se comprueba que ese echo retorne correctamente lo enviado.
Para levantar este contenedor temporal se usa `docker run --rm --platform linux/amd64 --network=tp0_testing_net` donde se indica:
* --rm se eliminara automaticamente al terminar
* --platform linux/amd64 especifica la arquitectura para que se pueda ejecutar desde cualquiera
* --network define a que red vamos a testear
Sobre ella se corre luego con alpine el netcat para ver el resultado. 


### Ej4
De aca en adelante, se ejecuta con `make docker-compose-up` y luego `make docker-compose-logs` para ver los logs. Para eliminar los contenedores se usa `make docker-compose-down`.

Para esto tuve que implementar un catch a la señal de Sigterm en ambos el cliente y el server.

Del lado del server se implementa con `signal.signal` donde definimos que funcion handleará que señal, en este caso SIGTERM. A partir de ello creamos un handler el cual ejecuta el cierre. En el handler en si se loggea el cierre, se setea en un bool de estado, se cierra el socket de aceptacion y el de todos los clientes.

Del lado del client en el go se crea un channel con `os.Signal` y luego `signal.Notify` donde se define que en ese chanel se recibira la señal. Luego se genera al comienzo del cliente la go rutine que ejecuta la function de handle, que queda bloqueada a espera de la notificacion. Al ser notificado cierra el skt y setea un bool de estado. 

En ambos casos el flujo se ejecuta sin ninguna diferencia respecto a no tenerlo, pero se define un plan de accion claro ante el SIGTERM para que ambos cierren sus recursos y terminen de forma ordenada. Esto es fundamental para sistemas escalables. 


## Parte 2: Repaso de Comunicaciones

Para mi sistema de comunicacion opte por comunicar siempre primero en un bloque de 4 bytes el tamaño del mensaje a recibir y luego el mensaje en si. De esta manera se logra una rapida y eficiente implementacion. Quien envia tan solo debe generar un mensaje que concatene los bytes (en Big Endian) del tamaño del mensaje con el mensaje, y quien recibe siempre lee los primero 4 y averigua cuanto mas leer. Esto permite una implementacion simple y robusta, ademas, no hay nunca trafico innecesario como si lo habria si los paquetes fueran de tamaño fijo. 

Para prevenir escrituras y lecturas cortas ambas iteran comparando el tamaño esperado a enviar vs el tamaño acumulado que fue retornando read o write segun corresponda. 

En cuanto a los sockets elegi TCP ya que en este caso es preferible tener asegurado que hayan llegado las bets a hacerlo rapido. Imaginemos un caso en que uno apuesta y gana y no llego la apuesta, automaticamente queda descartado UDP. Ademas, permite que nos abstraigamos lo mas posible de toda la logica de envio.

### Ej5
Para la comunicacion de las bets use lo mas simple posible, unir los campos separandolos con un delimitador, en este caso `|`. La ventaja de esto es que es muy simple y permite correr un split para obtener del otro lado todos los campos de las bets. Como mencione arriba esto se complementaba con el tamaño previamente enviado. Quedaria algo asi el paquete a codificar {4B tamaño}{ID}|{NOMBRE}|{APELLIDO}|{DOCUMENTO}|NACIMIENTO}|{NUMERO}. 

Del lado del server al obtener un mensaje se obtienen primero 4 bytes con los cuales sabe cuanto medira el string que contiene "{ID}|{NOMBRE}|{APELLIDO}|{DOCUMENTO}|NACIMIENTO}|{NUMERO}" y con split('|') ya obtuviste todo. 
El servidor en esta etapa continua teniendo conexiones volatiles y lo mismo el cliente, quien solo comunica su bet, recibe confirmacion y termina su programa. Del lado del servidor termina su comunicacion con ese cliente pero mantiene abierto por si otro cliente se conecta.

### Ej6
En esta etapa el nuevo escalon es la modalidad de los batches. Esto implica que deja de ser una comunicacion tan lineal y pasa a ser mas iterativa. 
Para ello implemente primero el cliente. En el cliente primero se leen las apuestas de los archivos definidos en los volumes de docker y se agrupan en batches segun definidos por el config. Aca nuevamente toma valor que sean flexibles los paquetes ya que se optimiza al maximo cada envio segun el tamaño de batch definido. Del lado del cliente se empaqueta todo y se separa cada bet con un `\n`. De esta manera a ojos del protocolo es un largo string que envia nuevamente en formato 4bytes tamaño y luego el mensaje. Para finalizar el envio de los batches, envia otro mensaje indicando que su agencia termino de enviar los paquetes. 

Del lado del server se reciben los paquetes y se hacen dos splits, primero por `\n` obteniendo las bets individuales y luego por `|` para obtener los campos. El servidor sigue recibiendo respuestas de un mismo cliente hasta que obtiene la señal de paro donde entiende que este cliente ya termino su trabajo y deja de esperar nuevos paquetes, terminando la comunciacion. Nuevamente el servidor queda abierto a espera de nuevas posibles conexiones. 

### Ej7
Para esta variacion tuve que modificar ambos lados. 

Para este paso es necesario ir y volver entre diferentes conexiones, ya que el server y el cliente no persisten las mismas, por lo que opte por crear una nueva abstraccion del lado del servidor, las agency. Estas son encargadas de manejar la conexion con una agencia. Saben entenderlas y traducir al servidor su necesidad, ademas, le mantienen el estado. Esto ultimo toma mucha relevancia para poder coordinar el sorteo. Por cada nueva conexion el cliente se presenta, si es nuevo se crea una agency, si es preexistente se actualiza la conexion. Lo que yo hice fue crear una pseudo barrera en la que se espera que todas las agencies esten listas. 
Cuando un cliente termino de enviar sus apuestas y quiere obtener los ganadores envia una solicitud y la agency guarda este estado y las solicita. 
Si todas las agencys ya estan en estado de obtener los ganadores se avanza con el sorteo, sino se retorna un no al cliente, recordando la agency en que estado esta. Del lado del cliente se espera un momento y se intenta nuevamente. 
Eventualmente todas estan listas, el servidor corre el sorteo y le pasa la lista de ganadores a la agency. La agency filtra a solo su numero de agencia y se la comunica al cliente. Es importante notar que la agency es un servicio al server.

## Parte 3: Repaso de Concurrencia
Primero que nada opte por usar procesos en lugar de threading. El link provisto por la catedra habla de que "Two threads calling a function may take twice as much time as a single thread calling the function twice. The GIL can cause I/O-bound threads to be scheduled ahead of CPU-bound threads, and it prevents signals from being delivered.". Eso es directamente un liquidar a la implementacion. Es un trabajo intenso de I/O y usamos señales para handle el cierre del mismo. Debido a esto multithreading queda descartado y usamos multiprocesing, que permite un manejo concurrente del servidor. Particularmente permite que el socket aceptador siempre este vivo y que se lance un proceso por cada cliente que se conecta. 

El problema que esto trae es coordinar el acceso a los recursos compartidos, en este caso el archivo de bets y el momento del sorteo. Para el archivo de bets use un lock clasico y para el sorteo una barrera que espera a que todos lleguen antes de avanzar. 

### Ej8
Para esta etapa se solucionan parte de los problemas del 7. Principalmente deja de tener neceisdad una agency relacionada especificamente a cada conexion ya que las mismas pueden quedarse abiertas en cada proceso y esperar en una barrier clasica. Esto llevo
a que no haya mas una agencia por conexion, sino que simplemente la use como microservicio. El servidor le pasa a una agency una tarea para alguien y este la cumple y muere. Al no conservar ningun estado no es necesario que persista pero permite abstraer al server de las particularidades del sistema. 

Para el flujo del programa se mantiene similar. El servidor acepta a un cliente y lanza un proceso para que lo maneje. El proceso recibe las apuestas, toma el lock y las guarda, suelta el lock (el archivo es seccion critica). El cliente luego pide ver los ganadores y el proceso va en busca de ellos pero se encuentra la barrier. La misma gestiona la espera y eventualmente le permite pasar, donde toma el archivo para leer los ganadores. Envia a traves de la agency para filtrar los ganadores que aplican a esa agencia y termina la comunicacion con este cliente. Al finalizar el proceso se une al hilo principal terminando su ciclo de vida.


## Ejecución de los tests
Para la correcta ejecución de los mismos me fue necesario incluir el sleep previo al cierre de los clientes.
![image](https://github.com/user-attachments/assets/0579bf31-7c58-45c2-ba8f-cca6c1a90b75)
