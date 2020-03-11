package incoming

import (
	"rsps-comm-test/internal/game"
	"rsps-comm-test/internal/systems"
	"rsps-comm-test/pkg/packet"
)

type WindowStatusPacket struct {
	Mode   byte
	Width  uint16
	Height uint16
}

func (w *WindowStatusPacket) HandlePacket(player *game.Player,  packet *packet.Packet, systemList []systems.System) {

}
