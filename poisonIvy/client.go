// Special Thanks: https://github.com/dtgreene/ivy2
package poisonIvy

import (
	"bufio"
	"io"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type Client struct {
	port      io.ReadWriteCloser
	InboundQ  chan []byte
	OutboundQ chan []byte
	alive     bool
	mu        sync.Mutex
	timer     *time.Timer
}

func NewClient() *Client {
	return &Client{
		InboundQ:  make(chan []byte, 100),
		OutboundQ: make(chan []byte, 100),
	}
}

func (c *Client) Connect(comPort string) error {
	config := &serial.Config{Name: comPort, Baud: 115200}
	p, err := serial.OpenPort(config)
	if err != nil {
		return err
	}

	c.port = p
	c.alive = true

	go c.readLoop()
	go c.writeLoop()
	return nil
}

func (c *Client) readLoop() {
	reader := bufio.NewReader(c.port)
	for c.alive {
		buf := make([]byte, 4096)
		n, err := reader.Read(buf)
		if err == nil && n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			c.InboundQ <- data
		}
	}
}

func (c *Client) writeLoop() {
	for c.alive {
		msg := <-c.OutboundQ
		_, err := c.port.Write(msg)
		if err != nil {
			c.Disconnect()
		}

		time.Sleep(20 * time.Millisecond)
	}
}

func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alive = false
	if c.port != nil {
		c.port.Close()
	}
}

func (c *Client) KeepAlive() {
	msg := GetBaseMessage(257, false, false)
	c.OutboundQ <- msg
}
