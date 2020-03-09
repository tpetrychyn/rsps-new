package outgoing

import (
	"log"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

type PlayerUpdatePacket struct {
	Actor *models.Actor
}

func NewPlayerUpdatePacket(player *models.Actor) *PlayerUpdatePacket {
	return &PlayerUpdatePacket{Actor:player}
}

// [ [player pos] | [player up masks] ]
func (p *PlayerUpdatePacket) Build() []byte {
	buffer := utils.NewStream()

	for _, v := range p.Actor.NearbyPlayers.Get() {
		if v != p.Actor && (v == nil) {
			// todo remove player
			continue
		}
	}

	// teleport segment
	if p.Actor.Movement.LastPosition == nil {
		p.Actor.Movement.LastPosition = p.Actor.Movement.Position
	}
	diffX := p.Actor.Movement.Position.X - p.Actor.Movement.LastPosition.X
	diffZ := p.Actor.Movement.Position.Z - p.Actor.Movement.LastPosition.Z
	diffH := p.Actor.Movement.Position.Height - p.Actor.Movement.LastPosition.Height

	buffer.WriteBits(1,  1)
	buffer.WriteBits(1, 0) // no update pending
	buffer.WriteBits(2, 3)

	buffer.WriteBits(1, 0) // tiles within viewing distance
	buffer.WriteBits(2, uint(diffH & 0x3))
	buffer.WriteBits(5, uint(diffX & 0x1F))
	buffer.WriteBits(5, uint(diffZ & 0x1F))

	// player skip count segment
	//buffer.WriteBits(1, 0)
	//buffer.WriteBits(2 ,0)

	buffer.SkipByte()

	// external/players outside view region?
	// count 2045
	buffer.WriteBits(1, 0)
	buffer.WriteBits(2, 3)
	buffer.WriteBits(11, 2045)

	by := buffer.Flush()

	size := len(by)
	out := append([]byte{79, byte(size << 8), byte(size & 0xFF)}, by...)
	log.Printf("%+v", out)
	return out
}

// full skip example:
//target []byte{79, 0, 3, 0, 127, 244}
