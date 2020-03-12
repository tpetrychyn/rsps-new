package models

type DirectionType struct {
	OrientationValue int
	PlayerValue      int
	NpcValue         int
}

type dir struct {
	North     DirectionType
	NorthEast DirectionType
	East      DirectionType
	SouthEast DirectionType
	South     DirectionType
	SouthWest DirectionType
	West      DirectionType
	NorthWest DirectionType
	None      DirectionType
}

var Direction = &dir{
	None: DirectionType{
		OrientationValue: -1,
		PlayerValue:      -1,
		NpcValue:         -1,
	},
	NorthWest: DirectionType{
		OrientationValue: 0,
		PlayerValue:      5,
		NpcValue:         0,
	},
	North: DirectionType{
		OrientationValue: 1,
		PlayerValue:      6,
		NpcValue:         1,
	},
	NorthEast: DirectionType{
		OrientationValue: 2,
		PlayerValue:      7,
		NpcValue:         2,
	},
	West: DirectionType{
		OrientationValue: 3,
		PlayerValue:      3,
		NpcValue:         3,
	},
	East: DirectionType{
		OrientationValue: 4,
		PlayerValue:      4,
		NpcValue:         4,
	},
	SouthWest: DirectionType{
		OrientationValue: 5,
		PlayerValue:      0,
		NpcValue:         5,
	},
	South: DirectionType{
		OrientationValue: 6,
		PlayerValue:      1,
		NpcValue:         6,
	},
	SouthEast: DirectionType{
		OrientationValue: 7,
		PlayerValue:      2,
		NpcValue:         7,
	},
}

func DirectionFromDeltas(deltaX int, deltaY int) DirectionType {
	if deltaY == 1 {
		if deltaX == 1 {
			return Direction.NorthEast
		}
		if deltaX == 0 {
			return Direction.North
		}
		return Direction.NorthWest
	}
	if deltaY == -1 {
		if deltaX == 1 {
			return Direction.SouthEast
		}
		if deltaX == 0 {
			return Direction.South
		}
		return Direction.SouthWest
	}

	if deltaX == 1 {
		return Direction.East
	}
	if deltaX == -1 {
		return Direction.West
	}

	return Direction.None
}

func (d *DirectionType) GetDiagonalComponents() []DirectionType {
	if *d == Direction.NorthEast {
		return []DirectionType{Direction.North, Direction.East}
	}
	if *d == Direction.NorthWest {
		return []DirectionType{Direction.North, Direction.West}
	}
	if *d == Direction.SouthEast {
		return []DirectionType{Direction.South, Direction.East}
	}
	if *d == Direction.SouthWest {
		return []DirectionType{Direction.South, Direction.West}
	}

	return nil
}

func (d *DirectionType) IsDiagonal() bool {
	return *d == Direction.NorthEast || *d == Direction.NorthWest || *d == Direction.SouthEast || *d == Direction.SouthWest
}

func (d *DirectionType) GetDeltaX() int {
	switch *d {
	case Direction.SouthEast, Direction.NorthEast, Direction.East:
		return 1
	case Direction.SouthWest, Direction.NorthWest, Direction.West:
		return -1
	}
	return 0
}

func (d *DirectionType) GetDeltaY() int {
	switch *d {
	case Direction.NorthWest, Direction.NorthEast, Direction.North:
		return 1
	case Direction.SouthWest, Direction.SouthEast, Direction.South:
		return -1
	}
	return 0
}

func (d *DirectionType) GetOpposite() DirectionType {
	switch *d {
	case Direction.North:
		return Direction.South
	case Direction.South:
		return Direction.North
	case Direction.East:
		return Direction.West
	case Direction.West:
		return Direction.East
	case Direction.NorthWest:
		return Direction.SouthEast
	case Direction.NorthEast:
		return Direction.SouthWest
	case Direction.SouthEast:
		return Direction.NorthWest
	case Direction.SouthWest:
		return Direction.NorthEast
	}
	return Direction.None
}
