package collision

type objectTypeStruct struct {
	LengthwiseWall  int
	DiagonalWall    int
	FloorDecoration int
}

var ObjectTypes = &objectTypeStruct{
	LengthwiseWall:  0,
	DiagonalWall:    9,
	FloorDecoration: 22,
}
