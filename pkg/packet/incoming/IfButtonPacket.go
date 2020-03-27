package incoming

import (
	"encoding/binary"
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
	"log"
)

type IfButtonPacket struct {}

func (w *IfButtonPacket) HandlePacket(player *entities.Player,  packet *packet.Packet) {
	var hash int
	binary.Read(packet.Buffer, binary.BigEndian, &hash)

	var slot uint16
	binary.Read(packet.Buffer, binary.BigEndian, &slot)

	var item uint16
	binary.Read(packet.Buffer, binary.BigEndian, &item)

	log.Printf("hash %v slot %v item %v", hash, slot, item)
}

