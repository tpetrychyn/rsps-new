package incoming

import (
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/packet"
)

type MouseMovePacket struct {
}

func (w *MouseMovePacket) HandlePacket(player *game.Player,  packet *packet.Packet) {

}
