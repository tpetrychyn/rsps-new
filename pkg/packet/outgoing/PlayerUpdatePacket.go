package outgoing

import (
	"bytes"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
)

type PlayerUpdatePacket struct {
	bytes []byte
}

func NewPlayerUpdatePacket(p *models.Actor) *PlayerUpdatePacket {
	packet := &PlayerUpdatePacket{}
	stream := utils.NewStream()
	updateStream := utils.NewStream()

	for _, v := range p.NearbyPlayers.Get() {
		if v != p && (v == nil) {
			// todo remove player
			continue
		}
	}

	if p.Movement.Teleported {
		// teleport segment
		if p.Movement.LastPosition == nil {
			p.Movement.LastPosition = p.Movement.Position
		}
		diffX := p.Movement.Position.X - p.Movement.LastPosition.X
		diffZ := p.Movement.Position.Y - p.Movement.LastPosition.Y
		diffH := p.Movement.Position.Height - p.Movement.LastPosition.Height

		stream.WriteBits(1, 1)
		stream.WriteBits(1, p.UpdateMask.UpdateRequired()) // update pending
		stream.WriteBits(2, 3)

		stream.WriteBits(1, 0) // tiles within viewing distance
		stream.WriteBits(2, uint(diffH&0x3))
		stream.WriteBits(5, uint(diffX&0x1F))
		stream.WriteBits(5, uint(diffZ&0x1F))
	} else if p.Movement.WalkDirection != models.Direction.None {
		dx := DirectionDeltaX[p.Movement.WalkDirection.PlayerValue]
		dy := DirectionDeltaY[p.Movement.WalkDirection.PlayerValue]
		if !packet.playerRun(stream, p, &dx, &dy) {
			packet.playerWalk(stream, p, &dx, &dy)
		}
	} else {
		// player skip count segment
		stream.WriteBits(1, 0)
		stream.WriteBits(2, 0)
	}
	stream.SkipByte()

	// second step
	// dx 2 dy 0
	//running true
	// direction 8

	packet.appendUpdates(updateStream, p, false)

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

	packet.bytes = out
	return packet
}

// [ [player pos] | [player up masks] ]
func (p *PlayerUpdatePacket) Build() []byte {
	return p.bytes
}

func (p *PlayerUpdatePacket) playerRun(stream *utils.Stream, target *models.Actor, dx, dy *int) bool {
	if !(target.Movement.IsRunning && target.Movement.RunDirection != models.Direction.None) {
		return false
	}

	*dx += DirectionDeltaX[target.Movement.RunDirection.PlayerValue]
	*dy += DirectionDeltaY[target.Movement.RunDirection.PlayerValue]
	dir := getPlayerRunningDirection(*dx, *dy)
	if dir == -1 {
		return false
	}

	stream.WriteBits(1, 1)
	stream.WriteBits(1, 1) // update flag
	stream.WriteBits(2, 2)
	stream.WriteBits(4, uint(dir))

	return true
}

func (p *PlayerUpdatePacket) playerWalk(stream *utils.Stream, target *models.Actor, dx, dy *int) {
	stream.WriteBits(1, 1)
	stream.WriteBits(1, target.UpdateMask.UpdateRequired()) // update flag
	stream.WriteBits(2, 1)
	dir := getPlayerWalkingDirection(*dx, *dy)
	stream.WriteBits(3, uint(dir))
}

func (p *PlayerUpdatePacket) appendUpdates(updateStream *utils.Stream, target *models.Actor, updateAppearance bool) {
	if updateAppearance {
		target.UpdateMask.Appearance = true
	}

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
	if target.UpdateMask.Movement {
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

	if target.UpdateMask.Movement {
		if target.Movement.Teleported {
			updateStream.WriteByte(0xFF)
		} else if target.Movement.IsRunning {
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

var DirectionDeltaX = []int{-1, 0, 1, -1, 1, -1, 0, 1}
var DirectionDeltaY = []int{-1, -1, -1, 0, 0, 1, 1, 1}

func getPlayerWalkingDirection(dx, dy int) int {
	if dx == -1 && dy == -1 {
		return 0
	}
	if dx == 0 && dy == -1 {
		return 1
	}
	if dx == 1 && dy == -1 {
		return 2
	}
	if dx == -1 && dy == 0 {
		return 3
	}
	if dx == 1 && dy == 0 {
		return 4
	}
	if dx == -1 && dy == 1 {
		return 5
	}
	if dx == 0 && dy == 1 {
		return 6
	}
	if dx == 1 && dy == 1 {
		return 7
	}
	return -1
}

func getPlayerRunningDirection(dx, dy int) int {
	if dx == -2 && dy == -2 {
		return 0
	}
	if dx == -1 && dy == -2 {
		return 1
	}
	if dx == 0 && dy == -2 {
		return 2
	}
	if dx == 1 && dy == -2 {
		return 3
	}
	if dx == 2 && dy == -2 {
		return 4
	}
	if dx == -2 && dy == -1 {
		return 5
	}
	if dx == 2 && dy == -1 {
		return 6
	}
	if dx == -2 && dy == 0 {
		return 7
	}
	if dx == 2 && dy == 0 {
		return 8
	}
	if dx == -2 && dy == 1 {
		return 9
	}
	if dx == 2 && dy == 1 {
		return 10
	}
	if dx == -2 && dy == 2 {
		return 11
	}
	if dx == -1 && dy == 2 {
		return 12
	}
	if dx == 0 && dy == 2 {
		return 13
	}
	if dx == 1 && dy == 2 {
		return 14
	}
	if dx == 2 && dy == 2 {
		return 15
	}
	return -1
}