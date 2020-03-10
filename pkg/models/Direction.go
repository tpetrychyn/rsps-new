package models

type DirectionEnum = struct {
	PlayerValue int
	NpcValue    int
}

type dir struct {
	North     DirectionEnum
	NorthEast DirectionEnum
	East      DirectionEnum
	SouthEast DirectionEnum
	South     DirectionEnum
	SouthWest DirectionEnum
	West      DirectionEnum
	NorthWest DirectionEnum
	None      DirectionEnum
}

var Direction = &dir{
	NorthWest: DirectionEnum{
		PlayerValue: 5,
		NpcValue:    0,
	},
	North: DirectionEnum{
		PlayerValue: 6,
		NpcValue:    1,
	},
	NorthEast: DirectionEnum{
		PlayerValue: 7,
		NpcValue:    2,
	},
	West: DirectionEnum{
		PlayerValue: 3,
		NpcValue:    3,
	},
	East: DirectionEnum{
		PlayerValue: 4,
		NpcValue:    4,
	},
	SouthWest: DirectionEnum{
		PlayerValue: 0,
		NpcValue:    5,
	},
	South: DirectionEnum{
		PlayerValue: 1,
		NpcValue:    6,
	},
	SouthEast: DirectionEnum{
		PlayerValue: 2,
		NpcValue:    7,
	},
	None: DirectionEnum{
		PlayerValue: -1,
		NpcValue:    -1,
	},
}

func DirectionFromDeltas(deltaX int, deltaY int) DirectionEnum {
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
