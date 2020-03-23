package incoming

import (
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
)

type WindowStatusPacket struct {
	Mode   byte
	Width  uint16
	Height uint16
}

func (w *WindowStatusPacket) HandlePacket(player *entities.Player,  packet *packet.Packet) {

}
