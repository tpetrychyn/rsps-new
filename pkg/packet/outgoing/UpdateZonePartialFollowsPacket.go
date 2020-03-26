package outgoing

import "github.com/tpetrychyn/rsps-comm-test/pkg/models"

//    type: FIXED
//    opcode: 64
//    structure:
//      - name: x
//        type: BYTE
//        trans: NEGATE
//      - name: z
//        type: BYTE
//        trans: SUBTRACT

// Tells the client to refresh region objs?
type UpdateZonePartialFollowsPacket struct {
	Tile *models.Tile
}

func (u *UpdateZonePartialFollowsPacket) Build() []byte {
	return []byte{64, byte(-u.Tile.X), 128 - byte(u.Tile.Y)}
}
