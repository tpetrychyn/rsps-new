package models

type Movement struct {
	Direction    DirectionType
	Position     *Tile
	LastPosition *Tile
	IsRunning    bool
}

func NewMovement() *Movement {
	return &Movement{
		Position:     new(Tile),
		LastPosition: nil,
		Direction:    Direction.None,
	}
}
