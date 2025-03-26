package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmout int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	signalChannel chan os.Signal
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		signalChannel: make(chan os.Signal, 1),
	}

	signal.Notify(client.signalChannel, syscall.SIGTERM)
	
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if  err != nil {
		log.Errorf("action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return err
	}
	c.conn = conn

	clientPresentation := "clientid: " + c.config.ID
	send([]byte(clientPresentation), c.conn, c.config.ID)

	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	//ej4
	go c.ShutHandle()
	c.createClientSocket()
	
	bets, err := c.getBets()
	
	validateAction("apuesta_serializada", err != nil, err, c.config.ID)

	sentSize, err := c.sendBets(bets)

	validateAction("apuesta_serializada", err != nil || sentSize == 0, err, c.config.ID)
		
	_,_ = c.askWinners()
	
	_ = c.reciveWinners()
	
	c.conn.Close()

	// sleep de cara a los tests
	time.Sleep(5 * time.Second)
	
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) ShutHandle() {
	log.Infof("action: begin handle | result: success" )
	<-c.signalChannel
	c.conn.Close()
	log.Infof("action: end handle | result: success" )
}

func (c* Client) getBets() ([]Bet, error) {
	filePath := "./data/agency-" + string(c.config.ID) + ".csv"
	bets_file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("action: open_csv | result: fail | error: %v", err)
		return nil, err
	}
	defer bets_file.Close()

	var bets []Bet
	scanner := bufio.NewScanner(bets_file)
	for scanner.Scan() {
		line := scanner.Text()
		bet, err := betFromString(line, c.config.ID)
		if err != nil {
			log.Errorf("action: parse_bet | result: fail | error: %v", err)
			return nil, err
		}
		bets = append(bets, bet)
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("action: scan_file | result: fail | error: %v", err)
		return nil, err
	}

	if len(bets) == 0 {
		return nil, fmt.Errorf("no bets found in CSV")
	}

	return bets, nil
}

func (c* Client) sendBets(bets []Bet) (int, error) {
	totalSent := 0
	i := 0
	for i < len(bets) {
		data := []byte{}
		for j := 0; j < c.config.BatchMaxAmout && i < len(bets); j++ {
			data = append(data, []byte(bets[i].getBetSerialized())...)
			i++
		}
		
		sent, err := send(data, c.conn, c.config.ID)
		
		if err != nil {
			log.Errorf("action: batch send failed | result: fail | client_id: %v | error: %v",
			 	c.config.ID, err)
			return 0, err
		}
		totalSent += sent
	}

	finalizado := fmt.Sprintf("Agencia %s ha finalizado la carga", c.config.ID)

	send([]byte(finalizado), c.conn, c.config.ID)

	return totalSent, nil
}

func (c Client) askWinners()(int, error) {
	return send([]byte("WINNERS"), c.conn, c.config.ID)
}

func send(data []byte, connection net.Conn, id string) (int, error) {
	totalSent := 0
	var err error

	dataSize := uint32(len(data))

	buffer := new(bytes.Buffer)
	
	err = binary.Write(buffer, binary.BigEndian, dataSize)
	
	validateSend("buff msg size", err, id)
	
	err = binary.Write(buffer, binary.BigEndian, data)
	
	validateSend("buff msg", err, id)
	
	messageSize := buffer.Len()
	for totalSent < messageSize {
		n, _ := connection.Write(buffer.Bytes())
		totalSent += n
	}

	return totalSent, nil
}

func (c* Client) reciveWinners() bool {
	if c.conn == nil {
		log.Infof("action: recv size | result: fail | client_id: %v | error: nil conn",
		c.config.ID,
		)
		return false
	}
	msgSize := 0
	sizeBuffer := make([]byte, 4)
	_, err := io.ReadFull(c.conn, sizeBuffer) 
	validateRecv("recv msg size", err, c.config.ID)

	msgSize = int(binary.BigEndian.Uint32(sizeBuffer))
	msgBuffer := make([]byte, msgSize)
	_, err = io.ReadFull(c.conn, msgBuffer)
	validateRecv("recv msg", err, c.config.ID)

	parts := strings.Split(string(msgBuffer), ": ")
	if len(parts) != 2 {
		return false
	}
	cant_ganadores := parts[1]

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %s", cant_ganadores)

	return true
}
