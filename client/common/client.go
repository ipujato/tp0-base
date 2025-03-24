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

	// ej4
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
	running bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		signalChannel: make(chan os.Signal, 1),
		running: true,
	}

	signal.Notify(client.signalChannel, syscall.SIGTERM)
	
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	validateAction("connect", err != nil, err, c.config.ID)
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
	
	log.Infof("action: loop_starts | result: success | client_id: %v", c.config.ID)

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount && c.running; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		log.Infof("action: loop_iter | result: success | client_id: %v", c.config.ID)
		
		//ej4
		if !c.running {
			log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
			break
		}

		// ej5
		
		// recibir y serializar la apuesta
		bets, err := c.getBets()
		
		validateAction("apuesta_serializada", err != nil, err, c.config.ID)

		// enviar con cuidado de que cubra bien la cantidad
		sentSize, err := c.sendBets(bets)
		
		time.Sleep(c.config.LoopPeriod)

		validateAction("apuesta_serializada", err != nil || sentSize == 0, err, c.config.ID)

		// recibir respuesta con cuidado de recibir el tamaño de la respuesta
		// msg, err := c.recvBetConfirmation()

		// validateAction("confimacion recibida", err != nil || msg == "", err, c.config.ID)

		// log.Infof("confimacion recibida | result: succes | msg: %s", msg)

		// log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", bet.Documento, bet.Numero)
		msg, err := c.reciveWinners()

		validateAction("ganadores recibida", err != nil || msg == "", err, c.config.ID)

		log.Infof("ganadores recibida | result: succes | msg: %s", msg)
		
		c.running = false
		c.conn.Close()
		
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
		
	}
	
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) ShutHandle() {
	log.Infof("action: begin handle | result: success" )
	<-c.signalChannel
	c.conn.Close()
	c.running = false
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
		bet, err := c.parseBet(line)
		if err != nil {
			log.Errorf("action: parse_bet | result: fail | error: %v", err)
			continue
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

	log.Infof("leidas %v apuestas", len(bets))

	return bets, nil
}

func (c* Client) parseBet(line string) (Bet, error) {
	splitedString := strings.Split(line, ",")
	if len(splitedString) != 5 {
		return Bet{}, fmt.Errorf("invalid bet format")
	}
	bet := Bet {
		Agencia: c.config.ID,
		Nombre: splitedString[0],
		Apellido: splitedString[1],
		Documento: splitedString[2],
		Nacimiento: splitedString[3],
		Numero: splitedString[4],
	}
	return bet, nil
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
		// log.Infof("escritura de conn en sendBets")
		n, _ := connection.Write(buffer.Bytes())
		// validateSend("send buff", err, id)
		totalSent += n
	}

	return totalSent, nil
}

func (c* Client) recvBetConfirmation() (string, error) {
	if c.conn == nil {
		log.Infof("action: recv size | result: fail | client_id: %v | error: nil conn",
			c.config.ID,
		)
		return "", nil
	}

	sizeBuffer := make([]byte, 4)
	_, err := io.ReadFull(c.conn, sizeBuffer) 
	validateRecv("recv msg size", err, c.config.ID)

	msgSize := int(binary.BigEndian.Uint32(sizeBuffer))
	msgBuffer := make([]byte, msgSize)
	_, err = io.ReadFull(c.conn, msgBuffer)
	validateRecv("recv msg", err, c.config.ID)

	log.Infof(string(msgBuffer))


	return string(msgBuffer), nil
}

func (c* Client) reciveWinners() (string, error) {
	if c.conn == nil {
		log.Infof("action: recv size | result: fail | client_id: %v | error: nil conn",
		c.config.ID,
		)
		return "", nil
	}
	log.Infof("entre a winners")
	msgSize := 0
	msgBuffer := []byte{}
	for msgSize==0 {
		sizeBuffer := make([]byte, 4)
		_, err := io.ReadFull(c.conn, sizeBuffer) 
		validateRecv("recv msg size", err, c.config.ID)
	
		msgSize = int(binary.BigEndian.Uint32(sizeBuffer))
		msgBuffer = make([]byte, msgSize)
		_, err = io.ReadFull(c.conn, msgBuffer)
		validateRecv("recv msg", err, c.config.ID)
	}
	log.Infof("recibi %s de tamaño %v", string(msgBuffer), msgSize)

	parts := strings.Split(string(msgBuffer), ": ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid winners format")
	}
	cant_ganadores := parts[1]

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %s", cant_ganadores)
	

	return string(msgBuffer), nil
}
