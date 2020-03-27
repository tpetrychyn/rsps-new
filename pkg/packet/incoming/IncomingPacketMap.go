package incoming

import (
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet"
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
const MoveMinimapClickPacketId = 52
const EventAppletFocusPacketId = 73
const MapBuildCompletePacketId = 76
const MoveGameClickPacketId = 96

const IfButtonPacketId0 = 68
const IfButtonPacketId1 = 21
const IfButtonPacketId2 = 48
const IfButtonPacketId3 = 19

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
	MoveMinimapClickPacketId: {
		PacketType: VARIABLE_BYTE,
		Handler:    new(MoveGameClickPacket),
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
	IfButtonPacketId0: {
		PacketType: FIXED,
		Length:     8,
		Handler:    new(IfButtonPacket),
	},
	IfButtonPacketId1: {
		PacketType: FIXED,
		Length:     8,
		Handler:    new(IfButtonPacket),
	},
	IfButtonPacketId2: {
		PacketType: FIXED,
		Length:     8,
		Handler:    new(IfButtonPacket),
	},
	IfButtonPacketId3: {
		PacketType: FIXED,
		Length:     8,
		Handler:    new(IfButtonPacket),
	},
}
