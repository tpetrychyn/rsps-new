package incoming

import (
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/packet"
)

type IgnorePacket struct {}

func (i *IgnorePacket) HandlePacket(player *game.Player,  packet *packet.Packet) {

}
