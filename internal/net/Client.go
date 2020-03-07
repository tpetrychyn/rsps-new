package net

import (
	"github.com/gtank/isaac"
	"net"
	"rsps-comm-test/pkg/packet"
)

type Client struct {
	connection      net.Conn
	upstreamQueue   chan packet.UpstreamMessage
	downstreamQueue chan packet.DownstreamMessage
	encryptor        *isaac.ISAAC
	decryptor        *isaac.ISAAC
}

func NewClient(conn net.Conn, encryptor *isaac.ISAAC, decryptor *isaac.ISAAC) *Client {
	client := &Client{
		connection: conn,
		upstreamQueue: make(chan packet.UpstreamMessage, 64),
		downstreamQueue: make(chan packet.DownstreamMessage, 256),
		encryptor: encryptor,
		decryptor: decryptor,
	}

	go client.downstreamListener()
	return client
}

func (c *Client) EnqueueOutgoing(message packet.DownstreamMessage) {
	c.downstreamQueue <- message
}

func (c *Client) downstreamListener() {
	for {
		message := <- c.downstreamQueue
		byteArray := message.Build()
		byteArray[0] = byte(uint32(byteArray[0]) + (c.encryptor.Rand() & 0xFF))

		c.connection.Write(byteArray)
	}
}

func (c *Client) Close() {
	close(c.upstreamQueue)
	close(c.downstreamQueue)
}