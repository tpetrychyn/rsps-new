package outgoing

import "rsps-comm-test/pkg/models"

type PlayerUpdatePacket struct {
	Actor *models.Actor
}

func NewPlayerUpdatePacket(player *models.Actor) *PlayerUpdatePacket {
	return &PlayerUpdatePacket{Actor:player}
}
