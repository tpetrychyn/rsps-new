package outgoing

import (
	"bytes"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

type PlayerUpdatePacket struct {
	Actor *models.Actor
}

// [ [player pos] | [player up masks] ]
func (p *PlayerUpdatePacket) Build() []byte {
	stream := utils.NewStream()
	updateStream := utils.NewStream()

	for _, v := range p.Actor.NearbyPlayers.Get() {
		if v != p.Actor && (v == nil) {
			// todo remove player
			continue
		}
	}

	if p.Actor.UpdateMask.NeedsPlacement {
		// teleport segment
		if p.Actor.Movement.LastPosition == nil {
			p.Actor.Movement.LastPosition = p.Actor.Movement.Position
		}
		diffX := p.Actor.Movement.Position.X - p.Actor.Movement.LastPosition.X
		diffZ := p.Actor.Movement.Position.Y - p.Actor.Movement.LastPosition.Y
		diffH := p.Actor.Movement.Position.Height - p.Actor.Movement.LastPosition.Height

		stream.WriteBits(1, 1)
		stream.WriteBits(1, p.Actor.UpdateMask.UpdateRequired()) // update pending
		stream.WriteBits(2, 3)

		stream.WriteBits(1, 0) // tiles within viewing distance
		stream.WriteBits(2, uint(diffH&0x3))
		stream.WriteBits(5, uint(diffX&0x1F))
		stream.WriteBits(5, uint(diffZ&0x1F))
	} else if p.Actor.Movement.IsRunning {
		stream.WriteBits(1, 1)
		stream.WriteBits(1, p.Actor.UpdateMask.UpdateRequired()) // update flag
		stream.WriteBits(2, 2)
		stream.WriteBits(4, uint(p.Actor.Movement.Direction.PlayerValue))
	} else if p.Actor.Movement.Direction != models.Direction.None {
		stream.WriteBits(1, 1)
		stream.WriteBits(1, p.Actor.UpdateMask.UpdateRequired()) // update flag
		stream.WriteBits(2, 1)
		stream.WriteBits(3, uint(p.Actor.Movement.Direction.PlayerValue))
	} else {
		// player skip count segment
		stream.WriteBits(1, 0)
		stream.WriteBits(2, 0)
	}
	stream.SkipByte()

	p.appendUpdates(updateStream, p.Actor, false)

	// external/players outside view region?
	// count 2045
	stream.WriteBits(1, 0)
	stream.WriteBits(2, 3)
	stream.WriteBits(11, 2045)

	buffer := new(bytes.Buffer)
	updateBytes := updateStream.Flush()
	if len(updateBytes) > 1 {
		buffer.Write(stream.Flush())
		buffer.Write(updateBytes)
	} else {
		buffer.Write(stream.Flush())
	}

	by := buffer.Bytes()

	size := len(by)
	out := append([]byte{79, byte(size << 8), byte(size & 0xFF)}, by...)
	return out
}

func (p *PlayerUpdatePacket) appendUpdates(updateStream *utils.Stream, target *models.Actor, updateAppearance bool) {
	if updateAppearance {
		target.UpdateMask.Appearance = true
	}

	//if target.Movement.IsRunning {
	//	target.UpdateMask.NeedsPlacement = true
	//}

	if target.UpdateMask.UpdateRequired() == 0 {
		return
	}

	var excess = 0x8
	var mask int
	if target.UpdateMask.Hitmark {
		mask |= 0x40
	}
	if target.UpdateMask.Graphic {
		mask |= 0x200
	}
	if target.UpdateMask.NeedsPlacement {
		mask |= 0x1000
	}
	if target.UpdateMask.ForcedMovement {
		mask |= 0x400
	}
	if target.UpdateMask.ForcedChat {
		mask |= 0x20
	}
	if target.UpdateMask.FaceTile {
		mask |= 0x4
	}
	if target.UpdateMask.Appearance {
		mask |= 0x1
	}
	if target.UpdateMask.FaceActor {
		mask |= 0x2
	}
	if target.UpdateMask.PublicChat {
		mask |= 0x10
	}
	if target.UpdateMask.Animation {
		mask |= 0x80
	}

	if mask >= 0x100 {
		mask |= excess
		updateStream.WriteWordLE(uint(mask))
	} else {
		updateStream.WriteByte(byte(mask))
	}

	if target.UpdateMask.NeedsPlacement {
		if target.Movement.IsRunning {
			updateStream.WriteByte(2+128)
		} else {
			updateStream.WriteByte(1+128)
		}
	}

	if target.UpdateMask.Appearance {
		pu := &PlayerAppearance{Target: target}
		updateStream.Write(pu.Build())
	}
}
