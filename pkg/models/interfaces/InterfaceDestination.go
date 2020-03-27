package interfaces

type InterfaceDestinationType struct {
	InterfaceId       uint
	FixedChildId      int
	ResizeChildId     int
	ResizeListChildId int
}

type InterfaceDestinationName string

const (
	ChatBoxInterface InterfaceDestinationName = "ChatBoxInterface"
	UsernameInterface InterfaceDestinationName = "UsernameInterface"
	InventoryInterface InterfaceDestinationName = "InventoryInterface"
	LogoutInterface InterfaceDestinationName = "LogoutInterface"
)

var InterfaceDestinations = map[InterfaceDestinationName]*InterfaceDestinationType{
	ChatBoxInterface: {
		InterfaceId:       162,
		FixedChildId:      24,
		ResizeChildId:     29,
		ResizeListChildId: 31,
	},
	UsernameInterface: {
		InterfaceId:       163,
		FixedChildId:      19,
		ResizeChildId:     9,
		ResizeListChildId: 9,
	},
	InventoryInterface: {
		InterfaceId:       149,
		FixedChildId:      69,
		ResizeChildId:     71,
		ResizeListChildId: 71,
	},
	LogoutInterface: {
		InterfaceId:       182,
		FixedChildId:      76,
		ResizeChildId:     78,
		ResizeListChildId: 78,
	},
}
