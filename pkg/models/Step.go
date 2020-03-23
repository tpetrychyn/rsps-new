package models

type Step struct {
	*Tile
	Direction DirectionType
	Head      *Step
	Cost      int
}
