class Agency:
    def __init__(self, agency_num: int, connection):
        self.agency_num = agency_num
        self.connection = connection

    def __repr__(self):
        return f"Agency(agency_num={self.agency_num}, connection={self.connection})"