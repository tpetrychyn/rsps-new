package game

import (
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/archives"
	rsModels "github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models/collision"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
)

type World struct {
	Chunks  map[string]*models.Chunk
}

func (w *World) GetOrLoadChunk(tile *models.Tile) *models.Chunk {
	baseChunkTile := tile.ToChunkBase()
	chunkX, chunkY := baseChunkTile.X, baseChunkTile.Y
	chunkId := fmt.Sprintf("%d-%d", chunkX, chunkY)
	if chunk, ok := w.Chunks[chunkId]; ok {
		return chunk
	}

	chunk := &models.Chunk{
		Id:              chunkId,
		Coords:          baseChunkTile,
		CollisionMatrix: make([]*models.CollisionMatrix, 4),
	}
	w.Chunks[chunkId] = chunk

	regionId := tile.ToRegionId()
	landArchive := archives.NewLandLoader(utils.GetCache())
	objArray := landArchive.LoadObjects(regionId, utils.GlobalXteaDefs[uint16(regionId)])

	for _, obj := range objArray {
		if !chunk.Contains(&models.Tile{X: uint16(obj.WorldX), Y: uint16(obj.WorldY)}) {
			continue
		}

		def := utils.GetDefinitions().Objects[obj.Id]

		width, length := def.Width, def.Length
		if obj.Orientation == 1 || obj.Orientation == 3 {
			width = def.Length
			length = def.Width
		}

		if !unWalkable(def, int(obj.Type)) {
			continue
		}

		tile := &models.Tile{
			X:      uint16(obj.WorldX),
			Y:      uint16(obj.WorldY),
			Height: uint16(obj.Height),
		}
		if int(obj.Type) == collision.ObjectTypes.FloorDecoration {
			if def.Interactive && def.Solid {
				w.PutTile(tile, models.NESW...)
			}
		} else if int(obj.Type) >= collision.ObjectTypes.DiagonalWall && int(obj.Type) < collision.ObjectTypes.FloorDecoration {
			for dx := 0; dx < width; dx++ {
				for dy := 0; dy < length; dy++ {
					tile := &models.Tile{
						X:      uint16(obj.WorldX+dx),
						Y:      uint16(obj.WorldY+dy),
						Height: uint16(obj.Height),
					}
					w.PutTile(tile, models.NESW...)
				}
			}
		} else if int(obj.Type) == collision.ObjectTypes.LengthwiseWall {
			w.PutWall(tile, models.WNES[int(obj.Orientation)])
		} else if int(obj.Type) == collision.ObjectTypes.TriangularCorner || int(obj.Type) == collision.ObjectTypes.RectangularCorner {
			w.PutWall(tile, models.WNES_DIAGONAL[int(obj.Orientation)])
		} else if int(obj.Type) == collision.ObjectTypes.WallCorner {
			w.PutLargeCornerWall(tile, models.WNES_DIAGONAL[int(obj.Orientation)])
		}
	}

	mapArchive := archives.NewMapLoader(utils.GetCache())
	blocked, _ := mapArchive.LoadBlockedTiles(regionId)
	for _, v := range blocked {
		if !chunk.Contains(&models.Tile{X: uint16(v.X), Y: uint16(v.Y)}) {
			continue
		}
	}

	return chunk
}

func (w *World) PutTile(tile *models.Tile, dirs ...models.DirectionType) {
	chunk := w.GetOrLoadChunk(tile)
	if chunk.CollisionMatrix[tile.Height] == nil {
		chunk.CollisionMatrix[tile.Height] = models.NewCollisionMatrix(models.ChunkSize, models.ChunkSize)
	}
	chunk.CollisionMatrix[tile.Height].PutTile(int(tile.X), int(tile.Y), dirs...)
}

func (w *World) PutWall(tile *models.Tile, dir models.DirectionType) {
	w.PutTile(tile, dir)
	w.PutTile(tile.Step(dir), dir.GetOpposite())
}

func (w *World) PutLargeCornerWall(tile *models.Tile, dir models.DirectionType) {
	directions := dir.GetDiagonalComponents()
	w.PutTile(tile, directions...)
	for _, d := range directions {
		w.PutTile(tile.Step(d), d.GetOpposite())
	}
}

func unWalkable(def *rsModels.ObjectDef, typ int) bool {
	isSolidFloorDecoration := typ == collision.ObjectTypes.FloorDecoration && def.Interactive
	isRoof := typ >= collision.ObjectTypes.DiagonalInteractable && typ < collision.ObjectTypes.FloorDecoration && def.Solid
	isWall := (typ >= collision.ObjectTypes.LengthwiseWall && typ <= collision.ObjectTypes.RectangularCorner || typ == collision.ObjectTypes.DiagonalWall) && def.Solid
	isSolidInteractable := (typ == collision.ObjectTypes.DiagonalInteractable || typ == collision.ObjectTypes.Interactable) && def.Solid
	return isSolidFloorDecoration || isRoof || isWall || isSolidInteractable
}
