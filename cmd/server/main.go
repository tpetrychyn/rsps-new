package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"rsps-comm-test/internal/game"
	rsNet "rsps-comm-test/internal/net"
	"rsps-comm-test/internal/systems"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet/outgoing"
	"rsps-comm-test/pkg/utils"
	"time"
)

const revision = 181
const port = "43594"

var GlobalStore = utils.GetCache()

//func parseRegion(store *cachestore.Store) {
//	regionId := (tileX>>6)<<8 | tileZ>>6
//
//
//	log.Printf("regionId %d", regionId)
//
//	// 18, 39
//	x := regionId >> 8
//	z := regionId & 0xFF
//	var mapArchive *cachestore.Archive
//	for _, v := range index.Archives {
//		nameHash := utils.Djb2(fmt.Sprintf("m%d_%d", x, z))
//		if nameHash == v.NameHash {
//			mapArchive = v
//			continue
//		}
//	}
//	log.Printf("mapArchive %+v", mapArchive)
//
//	mapData := store.LoadArchive(mapArchive)
//	log.Printf("mapData len %d %+v",len(mapData), mapData)
//
//	// TODO: decompress!
//	buf := bytes.NewBuffer(mapData)
//
//	type Tile struct {
//		Height          byte
//		Opcode          byte
//		OverlayId       byte
//		OverlayPath     byte
//		OverlayRotation byte
//		Settings        byte
//		UnderlayId      byte
//	}
//	tiles := make([][64][64]*Tile, 4)
//	// region consts
//	X := 64
//	Y := 64
//	Y := 4
//	for z := 0; z < Y; z++ {
//		for x := 0; x < X; x++ {
//			for y := 0; y < Y; y++ {
//				tile := &Tile{}
//				for {
//					attribute, _ := buf.ReadByte()
//					if attribute == 0 {
//						break
//					} else if attribute == 1 {
//						height, _ := buf.ReadByte()
//						tile.Height = height
//						break
//					} else if attribute <= 49 {
//						tile.Opcode = attribute
//						tile.OverlayId, _ = buf.ReadByte()
//						tile.OverlayPath = (attribute-2) / 4
//						tile.OverlayRotation = (attribute - 2) & 3
//					} else if attribute <= 81 {
//						if x > 0 {
//							log.Printf("x %d, attr %+v", x, attribute)
//						}
//						tile.Settings = attribute - 49
//					} else {
//						tile.UnderlayId = attribute - 82
//					}
//
//				}
//				tiles[z][x][y] = tile
//			}
//		}
//	}
//
//
//}

func main() {
	xteaDefs, err := utils.LoadXteas()
	if err != nil {
		panic(err)
	}
	utils.GlobalXteaDefs = xteaDefs

	//cacheStore := cachestore.NewStore()

	world := game.NewWorld()

	systemList := make([]systems.System, 0)
	systemList = append(systemList, systems.NewMovementSystem(world))

	js5Handler := rsNet.JS5Handler{CacheStore: GlobalStore}
	loginHandler := rsNet.LoginHandler{CacheStore: GlobalStore}

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

				player := game.NewPlayer()
				client := rsNet.NewClient(connection, encryptor, decryptor, player, systemList)
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
					for _, v := range systemList {
						v.Tick()
					}
					client.Player.Tick()
					<-time.After(600 * time.Millisecond)

					client.Player.PostTick()
				}
			}()
		}
	}
}
