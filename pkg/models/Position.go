package models

type Position struct {
	X      uint16
	Z      uint16
	Height uint16
}

func (p *Position) To30BitInt() int32 {
	return int32(p.Z) & 0x3FFF | (int32(p.X) & 0x3FFF) << 14 | (int32(p.Height) & 0x3) << 28
}
