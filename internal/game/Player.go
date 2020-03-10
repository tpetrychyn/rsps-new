package game

import (
	"rsps-comm-test/internal/components"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
	"rsps-comm-test/pkg/packet/outgoing"
)

type Player struct {
	*models.Actor
	MovementComponent *components.MovementComponent
	OutgoingQueue     chan packet.DownstreamMessage
}

func NewPlayer() *Player {
	actor := models.NewActor()
	actor.UpdateMask.NeedsPlacement = true
	actor.UpdateMask.Appearance = true
	return &Player{
		Actor: actor,
		MovementComponent: components.NewMovementComponent(actor.Movement),
	}
}

func (p *Player) Tick() {
	p.MovementComponent.Tick()

	p.OutgoingQueue <- &outgoing.PlayerUpdatePacket{Actor: p.Actor}
}

func (p *Player) PostTick() {
	p.Actor.UpdateMask.Clear()
	p.Movement.Direction = models.Direction.None
}