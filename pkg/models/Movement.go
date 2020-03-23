package models

type Movement struct {
	Teleported          bool
	WalkDirection       DirectionType
	RunDirection        DirectionType
	Position            *Tile
	LastPosition        *Tile
	IsRunning           bool
	LastKnownRegionBase *Tile
}

func NewMovement() *Movement {
	return &Movement{
		Position:            new(Tile),
		LastPosition:        nil,
		WalkDirection:       Direction.None,
		RunDirection:        Direction.None,
		LastKnownRegionBase: nil,
	}
}
