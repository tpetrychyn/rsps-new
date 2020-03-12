package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"osrs-cache-parser/pkg/archives"
	"osrs-cache-parser/pkg/archives/definitions"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/internal/game/entities"
	rsNet "rsps-comm-test/internal/net"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/models/collision"
	"rsps-comm-test/pkg/packet/outgoing"
	"rsps-comm-test/pkg/utils"
	"time"
)

const revision = 181
const port = "43594"

func main() {
	xteaDefs, err := utils.LoadXteas()
	if err != nil {
		panic(err)
	}
	utils.GlobalXteaDefs = xteaDefs

	cacheStore := utils.GetCache()
	js5Handler := rsNet.JS5Handler{CacheStore: cacheStore}
	loginHandler := rsNet.LoginHandler{CacheStore: cacheStore}

	objArchive := definitions.NewObjectArchive(cacheStore)
	utils.GetDefinitions().SetObjects(objArchive.LoadObjectDefs())

	world := &game.World{Regions: make(map[int]*models.Region), Chunks: make(map[string]*models.Chunk)}

	//12342
	landArchive := archives.NewLandArchive(utils.GetCache())
	objArray := landArchive.LoadObjects(12342, utils.GlobalXteaDefs[12342])

	for _, v := range objArray {
		baseX, baseY := ((12342>>8)&0xFF)<<6, (12342&0xFF)<<6
		x, y := baseX+v.LocalX, baseY+v.LocalY

		chunk := world.GetOrLoadChunk(&models.Tile{X: uint16(x), Y: uint16(y)})
		if chunk.CollisionMatrix[v.Height] == nil {
			chunk.CollisionMatrix[v.Height] = models.NewCollisionMatrix(models.ChunkSize, models.ChunkSize)
		}

		if int(v.Type) >= collision.ObjectTypes.DiagonalWall && int(v.Type) <= collision.ObjectTypes.FloorDecoration {
			chunk.CollisionMatrix[v.Height].AddDirs(x%models.ChunkSize, y%models.ChunkSize, models.NESW)
		}
	}

	chunk := world.GetOrLoadChunk(&models.Tile{X: 3094, Y: 3497})
	log.Printf("%+v", chunk.CollisionMatrix[0])

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer listener.Close()

	log.Printf("Listening on port %s", port)

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		reader := bufio.NewReader(connection)

		var requestType byte
		binary.Read(reader, binary.BigEndian, &requestType)

		if requestType == 15 {
			var gameVersion int32
			binary.Read(reader, binary.BigEndian, &gameVersion)
			if gameVersion != revision {
				connection.Write([]byte{6}) // out of date
				continue
			}

			connection.Write([]byte{0})
			go func() {
				// continuously loop through reading js5 request bytes for this socket until error
				for js5Handler.HandleRequest(connection, reader) {
				}
				connection.Close()
			}()
		}

		if requestType == 14 {
			go func() {
				encryptor, decryptor := loginHandler.HandleRequest(connection, reader)
				if encryptor == nil || decryptor == nil {
					connection.Close()
					return
				}

				player := entities.NewPlayer(world)
				client := rsNet.NewClient(connection, encryptor, decryptor, player)
				player.OutgoingQueue = client.DownstreamQueue

				player.Actor.Movement.Position = &models.Tile{
					X: 3094,
					Y: 3497,
				}

				client.DownstreamQueue <- &outgoing.RebuildLoginPacket{Position: player.Actor.Movement.Position}

				client.DownstreamQueue <- &outgoing.IfOpenTopPacket{} // main screen interface?

				client.DownstreamQueue <- &outgoing.RebuildNormalPacket{Position: &models.Tile{
					X: player.Actor.Movement.Position.X >> 3,
					Y: player.Actor.Movement.Position.Y >> 3,
				}}

				for {
					client.Player.Tick()
					<-time.After(600 * time.Millisecond)
					client.Player.PostTick()
				}
			}()
		}
	}
}
