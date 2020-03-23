package incoming

import (
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
)

type IgnorePacket struct {}

func (i *IgnorePacket) HandlePacket(player *entities.Player,  packet *packet.Packet) {

}
