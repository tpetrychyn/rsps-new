package incoming

import (
	"log"
	"rsps-comm-test/internal/game/entities"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
)

type MoveGameClickPacket struct {
	Z            uint16
	X            uint16
	MovementType byte
}

func (m *MoveGameClickPacket) HandlePacket(player *entities.Player, packet *packet.Packet) {
	z := packet.ReadShortA()
	x := packet.ReadShortA()
	_ = packet.ReadByte()

	player.MovementComponent.MoveTo(&models.Tile{
		X: x,
		Y: z,
	})

	log.Printf("move to %d %d", x, z)
}
