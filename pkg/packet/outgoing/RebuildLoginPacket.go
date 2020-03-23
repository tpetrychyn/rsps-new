package outgoing

import (
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
)

type RebuildLoginPacket struct {
	Position *models.Tile
}

func (r *RebuildLoginPacket) Build() []byte {
	rebuildLoginBuffer := utils.NewStream()
	rebuildLoginBuffer.WriteBits(30, uint(r.Position.To30BitInt()))

	for i:=0;i<2047;i++ {
		if i != 0 {
			rebuildLoginBuffer.WriteBits(18, 0)
		}
	}

	rebuildLoginBuffer.WriteWordLEA(uint(r.Position.Y >> 3))
	rebuildLoginBuffer.WriteWordA(uint(r.Position.X >> 3))

	normalXteas := (&RebuildNormalPacket{Position: &models.Tile{X: r.Position.X >> 3, Y: r.Position.Y >> 3}}).GetRegionXteas()
	rebuildLoginBuffer.Write(normalXteas)

	by := rebuildLoginBuffer.Flush()

	out := utils.NewStream()
	out.WriteByte(0)
	out.WriteWord(uint(len(by)))
	out.Write(by)

	return out.Flush()
}
