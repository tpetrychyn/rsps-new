package outgoing

import "github.com/tpetrychyn/rsps-comm-test/pkg/utils"

//     type: FIXED
//    opcode: 84
//    structure:
//      - name: top
//        order: LITTLE
//        type: SHORT
//        trans: ADD
type IfOpenTopPacket struct {
	Top int
}

func (i *IfOpenTopPacket) Build() []byte {
	stream := utils.NewStream()

	stream.WriteByte(84)
	stream.WriteWordLEA(uint(i.Top))
	return stream.Flush()
}
