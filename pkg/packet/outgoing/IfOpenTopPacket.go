package outgoing

type IfOpenTopPacket struct {}

func (i *IfOpenTopPacket) Build() []byte {
	return []byte{84, 33, 0}
}