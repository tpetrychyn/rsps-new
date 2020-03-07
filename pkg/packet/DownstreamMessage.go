package packet

type DownstreamMessage interface {
	Build() []byte
}
