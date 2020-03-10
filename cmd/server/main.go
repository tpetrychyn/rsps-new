package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"osrs-cache-parser/pkg/cachestore"
	rsNet "rsps-comm-test/internal/net"
	"rsps-comm-test/pkg/utils"
	"time"
)

const revision = 181
const port = "43594"

func parseRegion(store *cachestore.Store) {
	index := store.FindIndex(5) // maps index
	chunkX := 380
	chunkZ := 431

	tileX := (chunkX + 6) << 3
	tileZ := (chunkZ + 6) << 3

	regionId := (tileX>>6)<<8 | tileZ>>6


	log.Printf("regionId %d", regionId)

	// 18, 39
	x := regionId >> 8
	z := regionId & 0xFF
	var mapArchive *cachestore.Archive
	for _, v := range index.Archives {
		nameHash := utils.Djb2(fmt.Sprintf("m%d_%d", x, z))
		if nameHash == v.NameHash {
			mapArchive = v
			continue
		}
	}
	log.Printf("mapArchive %+v", mapArchive)

	mapData := store.LoadArchive(mapArchive)
	log.Printf("mapData len %d %+v",len(mapData), mapData)

	// TODO: decompress!
	buf := bytes.NewBuffer(mapData)

	type Tile struct {
		Height          byte
		Opcode          byte
		OverlayId       byte
		OverlayPath     byte
		OverlayRotation byte
		Settings        byte
		UnderlayId      byte
	}
	tiles := make([][64][64]*Tile, 4)
	// region consts
	X := 64
	Y := 64
	Z := 4
	for z := 0; z < Z; z++ {
		for x := 0; x < X; x++ {
			for y := 0; y < Y; y++ {
				tile := &Tile{}
				for {
					attribute, _ := buf.ReadByte()
					if attribute == 0 {
						break
					} else if attribute == 1 {
						height, _ := buf.ReadByte()
						tile.Height = height
						break
					} else if attribute <= 49 {
						tile.Opcode = attribute
						tile.OverlayId, _ = buf.ReadByte()
						tile.OverlayPath = (attribute-2) / 4
						tile.OverlayRotation = (attribute - 2) & 3
					} else if attribute <= 81 {
						if x > 0 {
							log.Printf("x %d, attr %+v", x, attribute)
						}
						tile.Settings = attribute - 49
					} else {
						tile.UnderlayId = attribute - 82
					}

				}
				tiles[z][x][y] = tile
			}
		}
	}

	baseX := ((regionId >> 8) & 0xFF) << 6
	baseY := (regionId & 0xFF) << 6
	for height:=0;height<Z;height++ {
		for lx:=0;lx<X;lx++ {
			for lz:=0;lz<Y;lz++ {
				tile := tiles[height][lx][lz]
				if tile.Settings & 0x1 == 0x1 {
					log.Printf("blocked tile at %d %d %d %d", lx, lz, baseX + lx, baseY + lz)
				}
			}
		}
	}
	log.Printf("tile %+v", tiles[0][0][0])
}

func main() {
	xteaDefs, err := utils.LoadXteas()
	if err != nil {
		panic(err)
	}
	utils.GlobalXteaDefs = xteaDefs

	cacheStore := cachestore.NewStore()

	parseRegion(cacheStore)


	js5Handler := rsNet.JS5Handler{CacheStore: cacheStore}
	loginHandler := rsNet.LoginHandler{CacheStore: cacheStore}

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
				client := loginHandler.HandleRequest(connection, reader)
				if client == nil {
					connection.Close()
					return
				}

				for {
					client.Player.Tick()
					<-time.After(600 * time.Millisecond)

					client.Player.PostTick()
				}
			}()
		}
	}
}
