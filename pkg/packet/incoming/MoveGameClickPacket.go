package incoming

import (
	"log"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
)

type MoveGameClickPacket struct {
	Z            uint16
	X            uint16
	MovementType byte
}

func (m *MoveGameClickPacket) HandlePacket(player *game.Player, packet *packet.Packet) {
	z := packet.ReadShortA()
	x := packet.ReadShortA()
	_ = packet.ReadByte()

	player.MovementComponent.MoveTo(&models.Position{
		X:      x,
		Z:      z,
	})

	log.Printf("move to %d %d", x, z)
}
