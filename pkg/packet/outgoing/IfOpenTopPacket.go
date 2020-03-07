package outgoing

type IfOpenTopPacket struct {}

func (i *IfOpenTopPacket) Build() []byte {
	return []byte{33}
}