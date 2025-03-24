import socket

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
    
    def send(self, message):
        confirmation_size = len(message).to_bytes(4, byteorder="big")
        full_msg = confirmation_size + message
        
        self.sock.sendall(full_msg)