package models

type Tile struct {
	X      uint16
	Y      uint16
	Height uint16
}

func (p *Tile) To30BitInt() int32 {
	return int32(p.Y) & 0x3FFF | (int32(p.X) & 0x3FFF) << 14 | (int32(p.Height) & 0x3) << 28
}

func (p *Tile) ToRegionId() int {
	return int((p.X>>6)<<8 | p.Y>>6)
}