import os
import socket
import logging
import signal
import traceback
from .utils import *
from .agency import *
from multiprocessing import Process, Barrier, Lock


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self.running = True
         
        signal.signal(signal.SIGTERM, self.__handle_shutdown)
        
        
        self.processes = []
        self.bets_lock = Lock()
        self.barrier = Barrier(int(os.getenv('EXPECTED_CLIENTS', '0')))

    def run(self):
        while self.running:
            client_sock = self.__accept_new_connection()
            if client_sock != None:
                conn = Connection(client_sock)
                p = Process(target=self.__handle_client_connection, args=(conn,))
                self.processes.append(p)
                p.start()

        for process in self.processes:
            try: 
                process.join()
            except:
                logging.error("Process could not be joined")

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
            agency_num = message.split(':')[1]

            Agency(int(agency_num), client_connection, self).recieve_msg()

            self.send_winners(agency_num, client_connection)

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            traceback.print_exc()
        finally:
            client_connection.close()


    def __handle_shutdown(self,  signum, frame):
        logging.info('action: shutdown | result: in_progress')
        self.running = False
        self._server_socket.close()
        for process in self.processes:
            process.terminate()
            process.join()
            self.processes.remove(process)
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
        
        logging.info(f'action: apuesta_recibida | result: success | cantidad: {counter}')
        self.bets_lock.acquire()
        store_bets(bets)
        self.bets_lock.release()


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
        
    def send_winners(self, agency_num, client_connection):
        self.barrier.wait()  
        winners_copy = self.__get_winners()
        logging.info('action: sorteo | result: success')
        Agency(int(agency_num), client_connection, self).check_for_winners(winners_copy) 
                
    def __get_winners(self):
        winners = []
        self.bets_lock.acquire()
        bets = load_bets()
        for bet in bets:
            if has_won(bet):
                winners.append(bet)
        self.bets_lock.release()
        return winners
        