package models

const ChunkSize = 8

type Chunk struct {
	Id              string
	Coords          *Tile
	CollisionMatrix []*CollisionMatrix // a matrix for each height level (4)
}

func (c *Chunk) ToTile(x, y uint16) *Tile {
	return &Tile{X: (x + 6) << 3, Y: (y + 6) << 3}
}

func (c *Chunk) Contains(tile *Tile) bool {
	return tile.ToChunkBase().X == c.Coords.X && tile.ToChunkBase().Y == c.Coords.Y
}
