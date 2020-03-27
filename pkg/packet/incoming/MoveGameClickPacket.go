package incoming

import (
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
	"log"
)

type MoveGameClickPacket struct {}

func (m *MoveGameClickPacket) HandlePacket(player *entities.Player, packet *packet.Packet) {
	y := packet.ReadShortA()
	x := packet.ReadShortA()
	_ = packet.ReadByte()

	log.Printf("x %d y %d", x, y)
	player.MovementComponent.MoveTo(&models.Tile{
		X: int(x),
		Y: int(y),
	})
}
