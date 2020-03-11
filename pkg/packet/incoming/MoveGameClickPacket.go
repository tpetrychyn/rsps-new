package incoming

import (
	"log"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/internal/systems"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet"
)

type MoveGameClickPacket struct {
	Z            uint16
	X            uint16
	MovementType byte
}

func (m *MoveGameClickPacket) HandlePacket(player *game.Player, packet *packet.Packet, systemList []systems.System) {
	z := packet.ReadShortA()
	x := packet.ReadShortA()
	_ = packet.ReadByte()

	for _, system := range systemList {
		if sys, ok := system.(*systems.MovementSystem); ok {
			sys.Add(player.Id, player.Movement, &models.Tile{
				X: x,
				Y: z,
			})
		}
	}

	//player.MovementComponent.MoveTo(&models.Tile{
	//	X: x,
	//	Y: z,
	//})

	//player.World.MovementSystem.Add(player.Movement, &models.Tile{
	//	X: x,
	//	Y: z,
	//})

	log.Printf("move to %d %d", x, z)
}
