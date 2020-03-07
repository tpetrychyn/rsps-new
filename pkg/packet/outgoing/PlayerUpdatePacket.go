package outgoing

import "rsps-comm-test/pkg/models"

type PlayerUpdatePacket struct {
	Actor *models.Actor
}

func NewPlayerUpdatePacket(player *models.Actor) *PlayerUpdatePacket {
	return &PlayerUpdatePacket{Actor:player}
}

func (p *PlayerUpdatePacket) Build() []byte {
	return []byte{79, 0, 3, 0, 127, 244}
}
