package models

type CollisionMatrix struct {
	Length int
	Width  int
	matrix []uint16 // y*c.Width + x: flag
}

type CollisionFlag int

const (
	PAWN_NORTH_WEST CollisionFlag = iota + 1
	PAWN_NORTH
	PAWN_NORTH_EAST
	PAWN_EAST
	PAWN_SOUTH_EAST
	PAWN_SOUTH
	PAWN_SOUTH_WEST
	PAWN_WEST
)

var NESW = []DirectionType{Direction.North, Direction.East, Direction.South, Direction.West}
var PawnFlags = []CollisionFlag{PAWN_NORTH_WEST, PAWN_NORTH, PAWN_NORTH_EAST, PAWN_WEST, PAWN_EAST, PAWN_SOUTH_WEST, PAWN_SOUTH, PAWN_SOUTH_EAST}

func NewCollisionMatrix(width, length int) *CollisionMatrix {
	return &CollisionMatrix{
		Length: length,
		Width:  width,
		matrix: make([]uint16, length * width),
	}
}

func (c *CollisionMatrix) AddDirs(x, y int, dirs []DirectionType) {
	for _, dir := range dirs {
		c.AddFlag(x, y, PawnFlags[dir.OrientationValue])
	}
}

func (c *CollisionMatrix) AddFlag(x, y int, flag CollisionFlag) {
	// x 6, 1
	idx := y*c.Width + x
	// idx 14
	c.matrix[idx] = c.matrix[idx] | uint16(1 << flag)
	// idx 14 should be 4
}

func (c *CollisionMatrix) IsBlocked(x, y int, direction DirectionType) bool {
	x = x % ChunkSize
	y = y % ChunkSize

	switch direction {
	case Direction.NorthWest:
		return c.hasFlag(x, y, PawnFlags[Direction.NorthWest.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.North.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.West.OrientationValue])
	case Direction.North:
		return c.hasFlag(x, y, PawnFlags[Direction.North.OrientationValue])
	case Direction.NorthEast:
		return c.hasFlag(x, y, PawnFlags[Direction.NorthEast.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.North.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.East.OrientationValue])
	case Direction.East:
		return c.hasFlag(x, y, PawnFlags[Direction.East.OrientationValue])
	case Direction.SouthEast:
		return c.hasFlag(x, y, PawnFlags[Direction.SouthEast.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.South.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.East.OrientationValue])
	case Direction.South:
		return c.hasFlag(x, y, PawnFlags[Direction.South.OrientationValue])
	case Direction.SouthWest:
		return c.hasFlag(x, y, PawnFlags[Direction.SouthWest.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.South.OrientationValue]) || c.hasFlag(x, y, PawnFlags[Direction.West.OrientationValue])
	case Direction.West:
		return c.hasFlag(x, y, PawnFlags[Direction.West.OrientationValue])
	}

	return true
}

func (c *CollisionMatrix) hasFlag(x, y int, flag CollisionFlag) bool {
	// flag = pawn west, bit = 8, ordinal = 7
	// flag.getBitAsShort.toInt() 256 ??
	idx := y*c.Width + x
	return c.matrix[idx]&0xFFFF&(1<<flag) != 0
}