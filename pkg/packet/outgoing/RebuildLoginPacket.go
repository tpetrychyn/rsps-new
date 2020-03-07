package outgoing

import (
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

type RebuildLoginPacket struct {
	Position *models.Position
}

func (r *RebuildLoginPacket) Build() []byte {
	rebuildLoginBuffer := utils.NewStream()
	rebuildLoginBuffer.WriteBits(30, uint(r.Position.To30BitInt()))
	// TODO: perhaps only send the players positions on this players region
	// seems weird to send all players positions on login lol
	for i:=0;i<2047;i++ {
		if i != 0 {
			rebuildLoginBuffer.WriteBits(18, 0)
		}
	}
	rebuildLoginBuffer.CloseBitAccess()

	rebuildLoginBuffer.WriteWordLEA(uint(r.Position.Z >> 3))
	rebuildLoginBuffer.WriteWordA(uint(r.Position.X >> 3))

	normalXteas := (&RebuildNormalPacket{Position: &models.Position{X: r.Position.X >> 3, Z: r.Position.Z >> 3}}).GetRegionXteas()
	rebuildLoginBuffer.Write(normalXteas)

	by := rebuildLoginBuffer.Flush()

	out := utils.NewStream()
	out.WriteByte(0)
	out.WriteWord(uint(len(by)))
	out.Write(by)

	return out.Flush()
}
