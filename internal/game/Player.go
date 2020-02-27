package game

import (
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet/outgoing"
)

type Player struct {
	Actor *models.Actor
}

func NewPlayer() *Player {
	return &Player{Actor: &models.Actor{
		Position:      &models.Position{},
		NearbyPlayers: make([]*models.Actor, 0),
		NearbyNpcs:    make([]*models.Actor, 0),
		OutgoingQueue: make(chan interface{}, 255)},
	}
}

func (p *Player) AppendOutgoing() {
	p.Actor.OutgoingQueue <- &outgoing.send
	p.Actor.OutgoingQueue <- &outgoing.SendPositionPacket{Position:&models.Position{}}
}
