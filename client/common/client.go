package common

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
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
	c.validateAction("connect", err != nil, err)

	c.conn = conn
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
		bet, err := c.getBets()
		
		c.validateAction("apuesta_serializada", err != nil, err)

		// enviar con cuidado de que cubra bien la cantidad
		sentSize, err := c.sendBets(bet)

		c.validateAction("apuesta_serializada", err != nil || sentSize == 0, err)

		// recibir respuesta con cuidado de recibir el tamaño de la respuesta
		msg, err := c.recvBetConfirmation()

		c.validateAction("apuesta_serializada", err != nil || msg == "", err)

		log.Infof("confimacion recibida | result: succes | msg: %s", msg)

		log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", bet.Documento, bet.Numero)

		c.running = false
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

func (c *Client) validateAction(action string, condition bool, err error) error {
	if condition {
		log.Errorf("action: %s | result: fail | client_id: %v | error: %v",
			action, c.config.ID, err)
		return err
	}
	return nil
}

func (c *Client) validateSend(action string, err error) (int, error) {
	if err != nil {
		log.Errorf("action: %s | result: fail | client_id: %v | error: %v",
			action, c.config.ID, err)
		return 0, err
	}
	return 0, nil
}

func (c *Client) validateRecv(action string, err error) (string, error) {
	if err != nil {
		log.Errorf("action: %s | result: fail | client_id: %v | error: %v",
			action, c.config.ID, err)
		return "", err
	}
	return "", nil
}

func (c Client) getBets() (Bet, error) {
	nombre := os.Getenv("NOMBRE")
	apellido := os.Getenv("APELLIDO")
	documento := os.Getenv("DOCUMENTO")
	nacimiento := os.Getenv("NACIMIENTO")
	numero := os.Getenv("NUMERO")

	var bet = Bet{
		Agencia: c.config.ID,
		Nombre:     nombre,
		Apellido:   apellido,
		Documento:  documento,
		Nacimiento: nacimiento,
		Numero:     numero,
	}

	log.Infof("bet created bet: %v", bet)

	return bet, nil
}

func (c Client) sendBets(bet Bet) (int, error) {
	data := []byte(bet.getBetSerialized())
	
	dataSize := uint32(len(data))
	buffer := new(bytes.Buffer)
	
	err := binary.Write(buffer, binary.BigEndian, dataSize)
	
	c.validateSend("buff bet size", err)
	
	err = binary.Write(buffer, binary.BigEndian, data)
	
	c.validateSend("buff bet", err)
	
	messageSize := buffer.Len()
	totalSent := 0
	for totalSent < messageSize {
		log.Infof("escritura de conn en sendBets")
		n, err := c.conn.Write(buffer.Bytes())
		c.validateSend("send buff", err)
		totalSent += n
	}

	log.Infof("action: send_bet | result: success | client_id: %v | bytes_sent: %v",
		c.config.ID,
		totalSent,
	)

	return totalSent, nil
}

func (c Client) recvBetConfirmation() (string, error) {
	// Esperamos recibir action: apuesta_almacenada | result: success | dni: 11111111 | numero: 1111 

	if c.conn == nil {
		log.Infof("action: recv size | result: fail | client_id: %v | error: nil conn",
			c.config.ID,
		)
		return "", nil
	}

	sizeBuffer := make([]byte, 4)
	_, err := io.ReadFull(c.conn, sizeBuffer) 
	c.validateRecv("recv msg size", err)

	msgSize := int(binary.BigEndian.Uint32(sizeBuffer))
	msgBuffer := make([]byte, msgSize)
	_, err = io.ReadFull(c.conn, msgBuffer)
	c.validateRecv("recv msg", err)

	c.conn.Close()

	return string(msgBuffer), nil
}
