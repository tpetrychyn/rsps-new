package incoming

import (
	"rsps-comm-test/internal/game/entities"
	"rsps-comm-test/pkg/packet"
)

type IgnorePacket struct {}

func (i *IgnorePacket) HandlePacket(player *entities.Player,  packet *packet.Packet) {

}
