package outgoing

import "github.com/tpetrychyn/rsps-comm-test/pkg/utils"

//  - message: gg.rsmod.game.message.impl.IfOpenSubMessage
//    type: FIXED
//    opcode: 77
//    structure:
//      - name: type
//        trans: ADD
//        type: BYTE
//      - name: overlay
//        order: MIDDLE
//        type: INT
//      - name: component
//        order: LITTLE
//        type: SHORT
//        trans: ADD
type IfOpenSubPacket struct {
	InterfaceId uint
	Parent      int
	Child       int
}

func (i *IfOpenSubPacket) Build() []byte {
	overlay := (i.Parent << 16) | i.Child

	stream := utils.NewStream()
	stream.WriteByte(77)
	stream.WriteByte(1 + 128) // type

	stream.WriteByte(byte(overlay >> 8)) // overlay
	stream.WriteByte(byte(overlay))
	stream.WriteByte(byte(overlay >> 24))
	stream.WriteByte(byte(overlay >> 16))

	stream.WriteWordLEA(i.InterfaceId) // component/interfaceId

	return stream.Flush()
}
