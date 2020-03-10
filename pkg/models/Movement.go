package models

type Movement struct {
	Direction    DirectionEnum
	Position     *Position
	LastPosition *Position
	IsRunning    bool
}

func NewMovement() *Movement {
	return &Movement{
		Position:     new(Position),
		LastPosition: nil,
		Direction:    Direction.None,
	}
}
