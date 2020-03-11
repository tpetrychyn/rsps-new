package game

import (
	"osrs-cache-parser/pkg/archives"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

type World struct {
	Regions map[int]*models.Region
}

func NewWorld() *World {
	return &World{
		Regions:        make(map[int]*models.Region),
	}
}

func (w *World) GetRegion(id int) *models.Region {
	// return loaded region
	if r, ok := w.Regions[id]; ok {
		return r
	}

	// load region def from archive
	mapArchive := archives.NewMapArchive(utils.GetCache())
	blocked, bridge := mapArchive.LoadBlockedTiles(id)

	region := &models.Region{
		Id:           id,
		BlockedTiles: blocked,
		BridgeTiles:  bridge,
	}
	w.Regions[id] = region
	return region
}
