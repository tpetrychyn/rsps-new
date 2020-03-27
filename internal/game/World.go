package game

import (
	"fmt"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/archives"
	rsModels "github.com/tpetrychyn/osrs-cache-parser/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
)

type World struct {
	Chunks map[string]*models.Chunk
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

	mapArchive := archives.NewMapLoader(utils.GetCache())
	blocked, bridges := mapArchive.LoadBlockedTiles(regionId)
	for _, v := range blocked {
		if !chunk.Contains(&models.Tile{X: v.X, Y: v.Y}) {
			continue
		}
		w.PutTile(&models.Tile{X: v.X, Y: v.Y, Height: v.Height}, models.NESW...)
	}

objectCollisionLoop:
	for _, obj := range objArray {
		if !chunk.Contains(&models.Tile{X: obj.WorldX, Y: obj.WorldY}) {
			continue
		}

		def := utils.GetDefinitions().Objects[obj.Id]
		if !unWalkable(def, int(obj.Type)) {
			continue
		}

		height := obj.Height
		for _, v := range bridges {
			if v.X == obj.WorldX && v.Y == obj.WorldY {
				if v.Height == obj.Height+1 {
					continue objectCollisionLoop
				}
				height--
			}
		}

		width, length := def.Width, def.Length
		if obj.Orientation == 1 || obj.Orientation == 3 {
			width = def.Length
			length = def.Width
		}

		tile := &models.Tile{
			X:      obj.WorldX,
			Y:      obj.WorldY,
			Height: height,
		}
		if int(obj.Type) == models.ObjectTypes.FloorDecoration {
			if def.Interactive && def.Solid {
				w.PutTile(tile, models.NESW...)
			}
		} else if int(obj.Type) >= models.ObjectTypes.DiagonalWall && int(obj.Type) < models.ObjectTypes.FloorDecoration {
			for dx := 0; dx < width; dx++ {
				for dy := 0; dy < length; dy++ {
					tile := &models.Tile{
						X:      obj.WorldX + dx,
						Y:      obj.WorldY + dy,
						Height: height,
					}
					w.PutTile(tile, models.NESW...)
				}
			}
		} else if int(obj.Type) == models.ObjectTypes.LengthwiseWall {
			w.PutWall(tile, models.WNES[int(obj.Orientation)])
		} else if int(obj.Type) == models.ObjectTypes.TriangularCorner || int(obj.Type) == models.ObjectTypes.RectangularCorner {
			w.PutWall(tile, models.WNES_DIAGONAL[int(obj.Orientation)])
		} else if int(obj.Type) == models.ObjectTypes.WallCorner {
			w.PutLargeCornerWall(tile, models.WNES_DIAGONAL[int(obj.Orientation)])
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
	isSolidFloorDecoration := typ == models.ObjectTypes.FloorDecoration && def.Interactive
	isRoof := typ >= models.ObjectTypes.DiagonalInteractable && typ < models.ObjectTypes.FloorDecoration && def.Solid
	isWall := (typ >= models.ObjectTypes.LengthwiseWall && typ <= models.ObjectTypes.RectangularCorner || typ == models.ObjectTypes.DiagonalWall) && def.Solid
	isSolidInteractable := (typ == models.ObjectTypes.DiagonalInteractable || typ == models.ObjectTypes.Interactable) && def.Solid
	return isSolidFloorDecoration || isRoof || isWall || isSolidInteractable
}
