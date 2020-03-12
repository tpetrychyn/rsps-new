package incoming

import (
	"rsps-comm-test/internal/game/entities"
	"rsps-comm-test/pkg/packet"
)

type MouseMovePacket struct {
}

func (w *MouseMovePacket) HandlePacket(player *entities.Player,  packet *packet.Packet) {

}
