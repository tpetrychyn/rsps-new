package models

const ChunkSize = 8

type Chunk struct {
	Id              string
	BlockedTiles    map[string]bool
	CollisionMatrix []*CollisionMatrix // a matrix for each height level (4)
}

func (c *Chunk) ToTile(x, y uint16) *Tile {
	return &Tile{X: (x + 6) << 3, Y: (y + 6) << 3}
}
