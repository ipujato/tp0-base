import socket
import logging
import signal
import struct
from .utils import *


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        # ej4
        self.clients = []
        self.running = True
        signal.signal(signal.SIGTERM, self.__handle_shutdown)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self.running:
            client_sock = self.__accept_new_connection()
            if client_sock != None:
                self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # while self.running:
                # # TODO: Modify the receive to avoid short-reads
                # msg = client_sock.recv(1024).rstrip().decode('utf-8')
                # ! recv
                ## rcv size
                esp_siz = b''
                read_size = 0
                while read_size < 4:
                        recvd = client_sock.recv(4 - read_size)
                        if not recvd:
                            raise Exception(f'Full read could not be achieved. Read up to now {read_size} of 4')
                        esp_siz += recvd
                        read_size += len(recvd)
                print(esp_siz)
                expected_size = struct.unpack('>I', esp_siz)[0]
                print(expected_size)

                expected_size = int(expected_size)
                if expected_size <= 0:
                    logging.error("action: receive_message | result: fail | error: non-positive size received")
                    return

                ## rcv msg
                read_size = 0
                message = b''
                try:
                    while read_size < expected_size:
                        recvd = client_sock.recv(expected_size - read_size)
                        if not recvd:
                            raise Exception(f'Full read could not be achieved. Read up to now: {read_size} of {expected_size}')
                        message += recvd
                        read_size += len(recvd)
                except Exception as e:
                    logging.error(f"action: receive_message | result: fail | error: {e}")
                    return
                    
                # ! store bet
                
                self.new_bet_management(message.decode('utf-8'))

                # # TODO: Modify the send to avoid short-writes

                addr = client_sock.getpeername()
                self.clients.append(addr)
                logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {message}')
                logging.info(f'1 socket {client_sock}')
                # ! confirm

                confirmation_to_send = "Bet received successfully".encode('utf-8')
                logging.info(f'2 socket {client_sock}')

                confirmation_size = struct.pack('>I', len(confirmation_to_send))
                logging.info(f'3 socket {client_sock}')
                message = confirmation_size + confirmation_to_send
                logging.info(f'4 socket {client_sock}')
                bytes_sent = 0
                logging.info(f'5 socket {client_sock}')

                client_sock.sendall(message.encode())
                logging.info(f'6 socket {client_sock}')
                # try:
                #     while bytes_sent < len(message):
                #         sent = client_sock.sendall(message[bytes_sent:])
                #         if sent == 0:
                #             raise RuntimeError("socket connection broken")
                #         bytes_sent += sent
                # except Exception as e:
                #     logging.error(f"action: send_confirmation | result: fail | error: {e}")

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            logging.info("server cierra el socket")
            client_sock.close()

    def __handle_shutdown(self,  signum, frame):
        logging.info('action: shutdown | result: in_progress')
        self.running = False
        self._server_socket.close()
        for client in self.clients:
            client.close()
        logging.info('action: shutdown | result: success')
        
    def new_bet_management(self, rcvd_bets: str):
        print("bet management")
        print(rcvd_bets)
        try:
            agencia, nombre, apellido, documento, nacimiento, numero = rcvd_bets.split('|')
            logging.info(f'action: new_bet_management | result: success | agencia: {agencia} | nombre: {nombre} | apellido: {apellido} | documento: {documento} | nacimiento: {nacimiento} | numero: {numero}')

            bets = []
            bet = Bet(agencia, nombre, apellido, documento, nacimiento, numero)
            bets.append(bet)
            store_bets(bets)
            logging.info(f'action: apuesta_almacenada | result: success | dni: {documento} | numero: {numero}')

        except Exception as e:
            logging.error(f"action: new_bet_management | result: fail | error: {e} | {rcvd_bets}")

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # ej4
        try:
            # Connection arrived
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except:
            logging.info('El socket se encuentra cerrado')
            return None