package models

import "math"

type Tile struct {
	X      uint16
	Y      uint16
	Height uint16
}

func (p *Tile) To30BitInt() int32 {
	return int32(p.Y) & 0x3FFF | (int32(p.X) & 0x3FFF) << 14 | (int32(p.Height) & 0x3) << 28
}

func (p *Tile) ToChunkBase() *Tile {
	return &Tile{X: (p.X >> 3) - 6, Y: (p.Y >> 3) - 6}
}

func (p *Tile) ToRegionId() int {
	return int((p.X >> 6) << 8 | p.Y >> 6)
}

func (p *Tile) Step(dir DirectionType) *Tile {
	return &Tile{
		X:      p.X + uint16(dir.GetDeltaX()),
		Y:      p.Y + uint16(dir.GetDeltaY()),
		Height: p.Height,
	}
}

func (p *Tile) IsWithinRadius(tile *Tile, radius int) bool {
	dx := int(math.Abs(float64(int(p.X) - int(tile.X))))
	dy := int(math.Abs(float64(int(p.Y) - int(tile.Y))))
	return dx <= radius && dy <= radius
}