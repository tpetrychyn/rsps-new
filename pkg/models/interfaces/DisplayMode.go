package interfaces

type DisplayModeType int

type displayMode struct {
	Fixed           DisplayModeType
	ResizableNormal DisplayModeType
	ResizableList   DisplayModeType
	Mobile          DisplayModeType
	Fullscreen      DisplayModeType
}

var DisplayMode = displayMode{
	Fixed:           0,
	ResizableNormal: 1,
	ResizableList:   2,
	Mobile:          3,
	Fullscreen:      4,
}

func (d DisplayModeType) GetDisplayComponentId() int {
	switch d {
	case DisplayMode.Fixed:
		return 548
	case DisplayMode.ResizableNormal:
		return 161
	case DisplayMode.ResizableList:
		return 164
	case DisplayMode.Fullscreen:
		return 165
	}
	return 548
}

func (d DisplayModeType) GetChildId() int {
	switch d {
	case DisplayMode.Fixed:
		return 548
	case DisplayMode.ResizableNormal:
		return 161
	case DisplayMode.ResizableList:
		return 164
	case DisplayMode.Fullscreen:
		return 165
	}
	return 548
}