package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/gtank/isaac"
	"net"
	"rsps-comm-test/internal/game/entities"
	"rsps-comm-test/pkg/packet"
	"rsps-comm-test/pkg/packet/incoming"
)

type Client struct {
	connection      net.Conn
	upstreamQueue   chan packet.UpstreamMessage
	DownstreamQueue chan packet.DownstreamMessage
	encryptor       *isaac.ISAAC
	decryptor       *isaac.ISAAC

	Player *entities.Player
}

func NewClient(conn net.Conn, encryptor *isaac.ISAAC, decryptor *isaac.ISAAC, player *entities.Player) *Client {
	client := &Client{
		connection:      conn,
		upstreamQueue:   make(chan packet.UpstreamMessage, 64),
		DownstreamQueue: make(chan packet.DownstreamMessage, 256),
		encryptor:       encryptor,
		decryptor:       decryptor,
		Player:          player,
	}

	go client.downstreamListener()
	go client.upstreamListener()
	return client
}

func (c *Client) downstreamListener() {
	for {
		message := <-c.DownstreamQueue
		byteArray := message.Build()
		byteArray[0] = byte(uint32(byteArray[0]) + (c.encryptor.Rand() & 0xFF))

		c.connection.Write(byteArray)
	}
}

func (c *Client) upstreamListener() {
	reader := bufio.NewReader(c.connection)
	for {
		by, err := reader.ReadByte()
		if err != nil {
			c.Close()
			return
		}
		opcode := byte(uint32(by) - (c.decryptor.Rand() & 0xFF))
		//log.Printf("opcode %+v", opcode)

		// map opcode to packet def'n
		packetDef := incoming.Packets[opcode]
		if packetDef == nil {
			continue
		}

		// find length of the packet
		var length uint16
		var payload []byte
		switch packetDef.PacketType {
		case incoming.VARIABLE_BYTE:
			byteLength, _ := reader.ReadByte()
			length = uint16(byteLength)
		case incoming.VARIABLE_SHORT:
			binary.Read(reader, binary.BigEndian, &length)
		case incoming.FIXED:
			length = packetDef.Length
		}

		// read the payload based on length
		payload = make([]byte, length)
		binary.Read(reader, binary.BigEndian, &payload)

		packetDef.Handler.HandlePacket(c.Player, &packet.Packet{
			Opcode:  opcode,
			Size:    length,
			Payload: payload,
			Buffer:  bytes.NewBuffer(payload),
		})
	}
}

func (c *Client) Close() {
	close(c.upstreamQueue)
	close(c.DownstreamQueue)
}
