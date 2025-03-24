import os
import socket
import logging
import signal
from .utils import *
from .agency import *
import time


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        # ej4
        self.agencies = []
        self.running = True
        self.expected_clients = int(os.getenv('EXPECTED_CLIENTS', '0'))
        logging.info(self.expected_clients)
        self.winners = []
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
                self.__handle_client_connection(Connection(client_sock))

    def __handle_client_connection(self, client_connection):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            esp_siz = client_connection.recieve_fixed_size_message(4)
            expected_size = int.from_bytes(esp_siz, byteorder="big")

            if expected_size <= 0:
                logging.error("action: receive_batch | result: fail | error: non-positive size received")
                return False
            
            message = client_connection.recieve_fixed_size_message(expected_size).decode('utf-8')
            
            position = self.__add_agency(message.split(':')[1], client_connection)

            self.agencies[position].recieve_bets()


            # amount = [0]
            # result = self.__recive_batches(client_sock, amount)
                            

            # if result:
            #     logging.info(f'action: apuestas totales para cliente | result: success | cantidad: {amount[0]}')
            #     answer_to_send = f'{amount[0]} bets saved successfully'.encode('utf-8')
            # else: 
            #     logging.info(f'action: apuestas totales para cliente | result: fail | cantidad: {amount[0]}')
            #     answer_to_send = f'{amount[0]} bet saved unsuccessfully'.encode('utf-8')

            # confirmation_size = len(answer_to_send).to_bytes(4, byteorder="big")
            # message = confirmation_size + answer_to_send
            
            # client_sock.sendall(message)

            self.__send_winners()

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_connection.close()


    def __handle_shutdown(self,  signum, frame):
        logging.info('action: shutdown | result: in_progress')
        self.running = False
        self._server_socket.close()
        for client in self.agencies:
            client.close()
        logging.info('action: shutdown | result: success')
        
    def new_bet_management(self, rcvd_bets, amount):
        bets = []
        counter = 0
        for bet in rcvd_bets.split('\n'):
            bet = bet.strip()
            if bet:
                try:
                    agencia, nombre, apellido, documento, nacimiento, numero = bet.split('|')
                    bet = Bet(agencia, nombre, apellido, documento, nacimiento, numero)
                    bets.append(bet)
                    amount[0] += 1
                    counter += 1
                except Exception as e:
                    logging.info(f'action: apuesta_recibida | result: fail | cantidad: {amount[0]}')
                    logging.error(f"action: new_bet_management | result: fail | error: {e} | {bet}")
        
        # logging.info(f'action: apuesta_recibida | result: success | cantidad: {counter}')
        store_bets(bets)

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
        
    def __send_winners(self):
        if len(self.agencies) == self.expected_clients:
            ready = True
            for agency in self.agencies:
                if not agency.is_ready():
                    ready = False
            
            if ready:
                self.__get_winners()
                logging.info('action: sorteo | result: success')
                for agency in self.agencies:
                    agency.check_for_winners(self.winners)
                
    def __get_winners(self):
        bets = load_bets()
        for bet in bets:
            if has_won(bet):
                self.winners.append(bet)
        logging.info(f'ganaron en total: {len(self.winners)} para {len(self.agencies)} de {self.expected_clients}')

    def __add_agency(self, agency_num, conn):
        self.agencies.append(Agency(int(agency_num), conn, self))
        return len(self.agencies)-1