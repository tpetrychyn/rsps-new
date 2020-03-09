package models

type Movement struct {
	Position     *Position
	LastPosition *Position
}

func NewMovement() *Movement {
	return &Movement{
		Position:     new(Position),
		LastPosition: nil,
	}
}