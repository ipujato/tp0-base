from threading import Lock
from .utils import *

class BetsMonitor:
    def __init__(self):
        self.lock = Lock()


    def store_bets(self, bets):
        with self.lock:
            store_bets(bets)

    def load_bets(self):
        with self.lock:
            return load_bets()
