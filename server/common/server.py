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
            amount = [0]
            result = self._recive_batches(client_sock, amount)
                            

            addr = client_sock.getpeername()
            self.clients.append(addr)
            if result:
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {amount[0]}')
                answer_to_send = "{amount[0]} bets saved successfully".encode('utf-8')
            else: 
                logging.info(f'action: apuesta_recibida | result: fail | cantidad: {amount[0]}')
                answer_to_send = "{amount[0]} bet saved unsuccessfully".encode('utf-8')

            confirmation_size = len(answer_to_send).to_bytes(4, byteorder="big")
            message = confirmation_size + answer_to_send
            
            client_sock.sendall(message)

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self.running = False

    def _recive_batches(self, client_sock, amount):
        still_receiving = True
        
        while still_receiving:
            esp_siz = self.__recieve_fixed_size_message(client_sock, 4)
            expected_size = int.from_bytes(esp_siz, byteorder="big")

            if expected_size <= 0:
                logging.error("action: receive_batch | result: fail | error: non-positive size received")
                return False

            try:
                message = self.__recieve_fixed_size_message(client_sock, expected_size)
            except Exception as e:
                logging.error(f"action: receive_batch | result: fail | error: {e}")
                return False
            if message.decode('utf-8') == "end":
                still_receiving = False
            else:
                try:
                    self.new_bet_management(message.decode('utf-8'), amount)
                except Exception as e:
                    logging.error(f"action: saving batch | result: fail | error: {e}")
                    return False
            
        return True

    def __recieve_fixed_size_message(self, client_sock, expected_size):
        read_size = 0
        esp_siz = b''
        while read_size < expected_size:
                recvd = client_sock.recv(expected_size - read_size)
                if not recvd:
                    raise Exception(f'Full read could not be achieved. Read up to now {read_size} of {expected_size}')
                esp_siz += recvd
                read_size += len(recvd)
        return esp_siz

    def __handle_shutdown(self,  signum, frame):
        logging.info('action: shutdown | result: in_progress')
        self.running = False
        self._server_socket.close()
        for client in self.clients:
            client.close()
        logging.info('action: shutdown | result: success')
        
    def new_bet_management(self, rcvd_bets, amount):
        # print("bet management")
        bets = []
        for bet in rcvd_bets.split('\n'):
            bet = bet.strip()
            if bet:
                try:
                    agencia, nombre, apellido, documento, nacimiento, numero = bet.split('|')
                    bet = Bet(agencia, nombre, apellido, documento, nacimiento, numero)
                    bets.append(bet)
                    amount[0] += 1
                except Exception as e:
                    logging.error(f"action: new_bet_management | result: fail | error: {e} | {bet}")

        store_bets(bets)
        logging.info('action: bets almacenadas')

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