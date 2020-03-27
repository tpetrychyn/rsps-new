package models

type objectTypeStruct struct {
	LengthwiseWall             int
	TriangularCorner           int
	WallCorner                 int
	RectangularCorner          int
	InteractableWallDecoration int
	InteractableWall           int
	DiagonalWall               int
	Interactable               int
	DiagonalInteractable       int
	FloorDecoration            int
}

var ObjectTypes = &objectTypeStruct{
	LengthwiseWall:             0,
	TriangularCorner:           1,
	WallCorner:                 2,
	RectangularCorner:          3,
	InteractableWallDecoration: 4,
	InteractableWall:           5,
	DiagonalWall:               9,
	Interactable:               10,
	DiagonalInteractable:       11,
	FloorDecoration:            22,
}
