from .utils import *
from .connection import *
import logging

class Agency:
    def __init__(self, agency_num: int, client_connection, loteria_nacional):
        self.agency_num = agency_num
        self.connection = client_connection
        self.loteria_nacional = loteria_nacional
        self.ready = False

    def close(self):
        self.connection.close()
    
    def recieve_bets(self):
        amount = [0]
        result = self.__recive_batches(amount)        

        if result:
            logging.info(f'action: apuestas totales para cliente | result: success | cantidad: {amount[0]}')
            answer_to_send = f'{amount[0]} bets saved successfully'.encode('utf-8')
        else: 
            logging.info(f'action: apuestas totales para cliente | result: fail | cantidad: {amount[0]}')
            answer_to_send = f'{amount[0]} bet saved unsuccessfully'.encode('utf-8')

        self.connection.send(answer_to_send)


    def __recive_batches(self, amount):
        still_receiving = True
        
        while still_receiving:
            esp_siz = self.connection.recieve_fixed_size_message(4)
            expected_size = int.from_bytes(esp_siz, byteorder="big")

            if expected_size <= 0:
                logging.error("action: receive_batch | result: fail | error: non-positive size received")
                return False

            try:
                message = self.connection.recieve_fixed_size_message(expected_size)
            except Exception as e:
                logging.error(f"action: receive_batch | result: fail | error: {e}")
                return False
            if message.decode('utf-8').__contains__("Agencia") and message.decode('utf-8').__contains__("ha finalizado la carga"):
                still_receiving = False
                self.ready = True

            else:
                try:
                    self.loteria_nacional.new_bet_management(message.decode('utf-8'), amount)
                except Exception as e:
                    logging.error(f"action: saving batch | result: fail | error: {e}")
                    return False
            
        return True
    
    def check_for_winners(self, winners):
        client_winners = [bet for bet in winners if bet.agency == self.agency_num]
        msg = "Ganaron:" + str(len(client_winners))
        logging.info(f"Para la agencia {self.agency_num} "+msg)
        self.connection.send(msg.encode('utf-8'))

    def is_ready(self):
        return self.ready

