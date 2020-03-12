package entities

import (
	"rsps-comm-test/internal/game"
	"rsps-comm-test/internal/game/components"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
	"rsps-comm-test/pkg/packet/outgoing"
)

type Player struct {
	*models.Actor
	World             *game.World
	MovementComponent *components.MovementComponent
	OutgoingQueue     chan packet.DownstreamMessage
}

func NewPlayer(world *game.World) *Player {
	actor := models.NewActor()
	actor.UpdateMask.NeedsPlacement = true
	actor.UpdateMask.Appearance = true
	return &Player{
		Actor:             actor,
		World:             world,
		MovementComponent: components.NewMovementComponent(actor.Movement, world),
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
