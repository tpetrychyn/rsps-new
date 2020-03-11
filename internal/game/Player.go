package game

import (
	"github.com/google/uuid"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
	"rsps-comm-test/pkg/packet/outgoing"
)

type Player struct {
	*models.Actor
	Id uuid.UUID
	//MovementComponent *systems.MovementSystem
	OutgoingQueue chan packet.DownstreamMessage
}

func NewPlayer() *Player {
	actor := models.NewActor()
	actor.UpdateMask.NeedsPlacement = true
	actor.UpdateMask.Appearance = true
	return &Player{
		Actor: actor,
		//MovementComponent: systems.NewMovementSystem(actor.Movement),
	}
}

func (p *Player) Tick() {
	//p.MovementComponent.Tick()

	p.OutgoingQueue <- &outgoing.PlayerUpdatePacket{Actor: p.Actor}
}

func (p *Player) PostTick() {
	p.Actor.UpdateMask.Clear()
	p.Movement.Direction = models.Direction.None
}
