package incoming

import (
	"log"
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
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

	log.Printf("move to %d %d, playerpos %d %d", x, z, player.Movement.Position.X, player.Movement.Position.Y)
}
