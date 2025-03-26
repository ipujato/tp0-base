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

type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmout int
}

type Client struct {
	config ClientConfig
	conn   net.Conn
	signalChannel chan os.Signal
	running bool
}

func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		signalChannel: make(chan os.Signal, 1),
		running: true,
	}

	signal.Notify(client.signalChannel, syscall.SIGTERM)
	
	return client
}

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

func (c *Client) StartClientLoop() {
	go c.ShutHandle()
	c.createClientSocket()
	
	log.Infof("action: loop_starts | result: success | client_id: %v", c.config.ID)

	for msgID := 1; msgID <= c.config.LoopAmount && c.running; msgID++ {
		log.Infof("action: loop_iter | result: success | client_id: %v", c.config.ID)
		if !c.running {
			log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
			return
		}
		
		bets, err := c.getBets()
		if err != nil {
			log.Errorf("action: obtencion de apuesta_serializadas | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			c.conn.Close()
			c.StartClientLoop()
		}

		sentSize, err := c.sendBets(bets)
		if err != nil || sentSize == 0 {
			log.Errorf("action: envio de apuesta_serializadas | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return
		}
		
		c.running = false
				
	}
	
	c.conn.Close()
	pending_results := true

	for pending_results {
		c.createClientSocket()
		_,_ = c.askWinners()
		success := c.reciveWinners()
		if success {
			pending_results = false
		}
		c.conn.Close()
	}

	// sleep por los tests
	time.Sleep(2 * time.Second) 
	
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
		bet, err := betFromString(line, c.config.ID)
		if err != nil {
			log.Errorf("action: parse_bet | result: fail | error: %v", err)
			return c.getBets()
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
	sent, err := send([]byte(finalizado), c.conn, c.config.ID)
	if err != nil {
		log.Errorf("action: batch end send failed | result: fail | client_id: %v | error: %v",
			 c.config.ID, err)
		return 0, err
	}
	totalSent += sent

	return totalSent, nil
}

func (c Client) askWinners()(int, error) {
	return send([]byte("WINNERS"), c.conn, c.config.ID)
}

func send(data []byte, connection net.Conn, id string) (int, error) {
	totalSent := 0
	dataSize := uint32(len(data))
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, dataSize)
	if err != nil {
		log.Errorf("action: wirte buff msg size | result: fail | client_id: %v | error: %v",
			id, err)
		return 0, err
	}
	
	err = binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Errorf("action: buff msg | result: fail | client_id: %v | error: %v",
			id, err)
		return 0, err
	}
	
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
	if err != nil && err != io.EOF{
		log.Errorf("action: recv size | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}

	msgSize = int(binary.BigEndian.Uint32(sizeBuffer))
	msgBuffer := make([]byte, msgSize)
	_, err = io.ReadFull(c.conn, msgBuffer)
	if err != nil && err != io.EOF{
		log.Errorf("action: recv winners | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}

	parts := strings.Split(string(msgBuffer), ": ")
	if len(parts) != 2 {
		return false
	}
	cant_ganadores := parts[1]

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %s", cant_ganadores)

	return true
}
