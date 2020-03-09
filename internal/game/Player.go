package game

import (
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
	"rsps-comm-test/pkg/packet/outgoing"
)

type Player struct {
	Actor *models.Actor
	OutgoingQueue chan packet.DownstreamMessage
}

func NewPlayer() *Player {
	return &Player{
		Actor: models.NewActor(),
	}
}

func (p *Player) Tick() {
	p.OutgoingQueue <- &outgoing.PlayerUpdatePacket{Actor: p.Actor}
}
