package outgoing

import "github.com/tpetrychyn/rsps-comm-test/pkg/models"

type MapFlagPacket struct {
	LastKnownRegionBase *models.Tile
	Tile *models.Tile
}

//    type: FIXED
//    opcode: 67
//    structure:
//      - name: x
//        type: BYTE
//      - name: z
//        type: BYTE
func (m *MapFlagPacket) Build() []byte {
	if m.LastKnownRegionBase == nil || m.Tile == nil {
		return []byte{67, 0xFF, 0xFF}
	}

	return []byte{67, byte(m.Tile.X-m.LastKnownRegionBase.X), byte(m.Tile.Y-m.LastKnownRegionBase.Y)}
}


