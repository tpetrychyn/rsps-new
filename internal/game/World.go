package game

import (
	"fmt"
	"log"
	"osrs-cache-parser/pkg/archives"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

type World struct {
	Regions map[int]*models.Region
	Chunks  map[string]*models.Chunk
}

func (w *World) GetOrLoadChunk(tile *models.Tile) *models.Chunk {
	chunkX, chunkY := tile.ToChunkCoords()
	chunkId := fmt.Sprintf("%d-%d", chunkX, chunkY)
	if chunk, ok := w.Chunks[chunkId]; ok {
		return chunk
	}

	// TODO: can I just make blocktiles set their collisionMatrix byte to 255?
	chunk := &models.Chunk{
		Id:              chunkId,
		BlockedTiles:    make(map[string]bool),
		CollisionMatrix: make([]*models.CollisionMatrix, 4),
	}

	w.Chunks[chunkId] = chunk
	return chunk
}

func (w *World) GetRegion(id int) *models.Region {
	if region, ok := w.Regions[id]; ok {
		return region
	}

	// load region def from archive
	mapArchive := archives.NewMapArchive(utils.GetCache())
	blocked, bridge := mapArchive.LoadBlockedTiles(id)

	landArchive := archives.NewLandArchive(utils.GetCache())
	// TODO: landarchive loading is broken
	worldObjects := landArchive.LoadObjects(id, utils.GlobalXteaDefs[uint16(id)])

	for k, v := range worldObjects {
		log.Printf("k %v v %v", k, v)
		//blocked[k] = true
	}

	region := &models.Region{
		Id:           id,
		BlockedTiles: blocked,
		BridgeTiles:  bridge,
	}
	x, y := region.GetBase()
	log.Printf("regionbase %v %v", x, y)
	w.Regions[id] = region
	return region
}
