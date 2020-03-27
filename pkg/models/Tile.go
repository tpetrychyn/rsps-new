package models

import "math"

type Tile struct {
	X      int
	Y      int
	Height int
}

func (p *Tile) To30BitInt() int32 {
	return int32(p.Y)&0x3FFF | (int32(p.X)&0x3FFF)<<14 | (int32(p.Height)&0x3)<<28
}

func (p *Tile) ToChunkBase() *Tile {
	return &Tile{X: (p.X >> 3) - 6, Y: (p.Y >> 3) - 6}
}

func (p *Tile) ToRegionId() int {
	return (p.X>>6)<<8 | p.Y>>6
}

func (p *Tile) ToLocal(other *Tile) *Tile {
	return &Tile{
		X:     ((other.X >> 3) - (p.X>>3))<<3,
		Y:      ((other.Y >> 3) - (int(p.Y)>>3))<<3,
		Height: p.Height,
	}
}

func (p *Tile) Step(dir DirectionType) *Tile {
	return &Tile{
		X:      p.X + dir.GetDeltaX(),
		Y:      p.Y + dir.GetDeltaY(),
		Height: p.Height,
	}
}

func (p *Tile) IsWithinRadius(tile *Tile, radius int) bool {
	dx := int(math.Abs(float64(p.X - tile.X)))
	dy := int(math.Abs(float64(p.Y - tile.Y)))
	return dx <= radius && dy <= radius
}

func (p *Tile) DistanceTo(other *Tile) int {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return int(math.Sqrt(float64(dx*dx) + float64(dy*dy)))
}