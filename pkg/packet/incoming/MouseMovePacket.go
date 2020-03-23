package incoming

import (
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
)

type MouseMovePacket struct {
}

func (w *MouseMovePacket) HandlePacket(player *entities.Player,  packet *packet.Packet) {

}
