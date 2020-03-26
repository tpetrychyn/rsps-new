package outgoing

import (
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
)

//      - name: tile
//        type: BYTE
//        trans: ADD
//      - name: settings
//        type: BYTE
//        trans: SUBTRACT
//      - name: id
//        type: SHORT
//        trans: ADD

// Spawns a world object
type LocAddChangePacket struct {
	Tile *models.Tile
	Type int
	Rot  int
	Id   int
}

func (l *LocAddChangePacket) Build() []byte {
	pos := ((int(l.Tile.X) & 0x7) << 4) | (int(l.Tile.Y) & 0x7)

	stream := utils.NewStream()
	stream.WriteByte(byte(pos + 128))
	stream.WriteByte(128 - byte((l.Type << 2) | l.Rot))
	stream.WriteWordA(uint(l.Id))

	by := stream.Flush()
	out := append([]byte{6}, by...)
	return out
}
