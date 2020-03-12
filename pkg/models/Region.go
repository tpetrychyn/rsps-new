package models

type Region struct {
	Id           int
	BlockedTiles map[string]bool
	BridgeTiles  map[string]bool
}

func CoordsToRegionId(x, y int) int {
	return (x>>6)<<8 | y>>6
}

func (r *Region) GetBase() (int, int) {
	return ((r.Id >> 8) & 0xFF) << 6, (r.Id & 0xFF) << 6
}