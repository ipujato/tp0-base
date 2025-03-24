import logging

class Connection:
    def __init__(self, client_sock):
        self.sock = client_sock

    def close(self):
        self.sock.close()

    def recieve_fixed_size_message(self, expected_size):
        read_size = 0
        esp_siz = b''
        while read_size < expected_size:
                recvd = self.sock.recv(expected_size - read_size)
                if not recvd:
                    raise Exception(f'Full read could not be achieved. Read up to now {read_size} of {expected_size}')
                esp_siz += recvd
                read_size += len(recvd)
        return esp_siz
    
    def send(self, encoded_message):
        if self.sock.fileno() == -1:  # El file descriptor -1 indica socket cerrado
            logging.error("action: send | result: fail | error: Trying to send on closed socket")
            return
        message_size = len(encoded_message).to_bytes(4, byteorder="big")
        full_msg = message_size + encoded_message
        
        self.sock.sendall(full_msg)