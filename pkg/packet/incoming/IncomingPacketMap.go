package incoming

import (
	"rsps-comm-test/internal/game/entities"
	"rsps-comm-test/pkg/packet"
)

type IncomingPacket interface {
	HandlePacket(player *entities.Player, packet *packet.Packet)
}

type PacketType int

const (
	FIXED PacketType = iota
	VARIABLE_BYTE
	VARIABLE_SHORT
)

type PacketDefinition struct {
	PacketType PacketType
	Length     uint16
	Handler    IncomingPacket
}

const IgnorePacketId = 4
const NoTimeoutPacketId = 22
const MouseMovePacketId = 34
const WindowStatusPacketId = 35
const MouseClickPacketId = 41
const EventAppletFocusPacketId = 73
const MapBuildCompletePacketId = 76
const MoveGameClickPacketId = 96

var Packets = map[byte]*PacketDefinition{
	IgnorePacketId: {
		PacketType: VARIABLE_BYTE,
		Handler:    new(IgnorePacket),
	},
	MouseMovePacketId: {
		PacketType: VARIABLE_BYTE,
		Handler:    new(MouseMovePacket),
	},
	WindowStatusPacketId: {
		PacketType: FIXED,
		Length:     5,
		Handler:    new(WindowStatusPacket),
	},
	MouseClickPacketId: {
		PacketType: FIXED,
		Length:     6,
		Handler:    new(IgnorePacket),
	},
	NoTimeoutPacketId: {
		PacketType: FIXED,
		Length:     0,
		Handler:    new(IgnorePacket),
	},
	EventAppletFocusPacketId: {
		PacketType: FIXED,
		Length:     1,
		Handler:    new(IgnorePacket),
	},
	MapBuildCompletePacketId: {
		PacketType: FIXED,
		Length:     0,
		Handler:    new(IgnorePacket),
	},
	MoveGameClickPacketId: {
		PacketType: VARIABLE_BYTE,
		Handler:    new(MoveGameClickPacket),
	},
}
