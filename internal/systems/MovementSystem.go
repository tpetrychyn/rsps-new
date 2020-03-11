package systems

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/models"
	"rsps/util"
)

type movementEntity struct {
	Movement *models.Movement
	Steps    []*models.Step
}

type MovementSystem struct {
	entities map[uuid.UUID]*movementEntity
	world    *game.World
	//Movement *models.Movement
	//Steps    []*models.Step
}

func NewMovementSystem(world *game.World) *MovementSystem {
	return &MovementSystem{
		entities: make(map[uuid.UUID]*movementEntity, 0),
		world:    world,
		//Movement: movement,
		//Steps:    make([]*models.Step, 0),
	}
}

func (m *MovementSystem) Tick() {
	for k, v := range m.entities {
		if len(v.Steps) == 0 {
			delete(m.entities, k)
			continue
		}

		walkPoint := v.Steps[0]
		m.entities[k].Steps = v.Steps[1:]

		log.Printf("vsteps %+v", v.Steps)

		if walkPoint != nil {
			v.Movement.Position = walkPoint.Tile
			v.Movement.Direction = walkPoint.Direction
			log.Printf("wdir %+v", walkPoint.Direction)
		}
	}

	//if len(m.Steps) == 0 {
	//	return
	//}
	//
	//walkPoint := m.Steps[0]
	//m.Steps = m.Steps[1:]
	//
	//var runPoint *models.Step
	//if m.Movement.IsRunning && len(m.Steps) > 0 {
	//	runPoint = m.Steps[0]
	//	m.Steps = m.Steps[1:]
	//}
	//
	//if walkPoint != nil {
	//	m.Movement.Position = walkPoint.Tile
	//	m.Movement.Direction = walkPoint.Direction
	//	log.Printf("wdir %+v", walkPoint.Direction)
	//}
	//
	//if runPoint != nil {
	//	m.Movement.Position = runPoint.Tile
	//	m.Movement.Direction = runPoint.Direction
	//	log.Printf("rdir %+v", runPoint.Direction)
	//}
}

func (m *MovementSystem) Add(id uuid.UUID, movement *models.Movement, destination *models.Tile) {
	entity := &movementEntity{
		Movement: movement,
		Steps:    make([]*models.Step, 0),
	}
	entity.addPosition(destination, m.world)
	m.entities[id] = entity
	//m.entities = append(m.entities, entity)
}

//func (m *MovementSystem) MoveTo(movement *models.Movement, p *models.Tile) {
//	m.Steps = make([]*models.Step, 0)
//	m.addPosition(p)
//}

func (e *movementEntity) addPosition(p *models.Tile, world *game.World) {
	last := e.getLast()
	x := int(p.X)
	y := int(p.Y)

	deltaX := x - int(last.Tile.X)
	deltaY := y - int(last.Tile.Y)

	max := util.Abs(deltaX)
	if util.Abs(deltaY) > util.Abs(deltaX) {
		max = util.Abs(deltaY)
	}

	for i := 0; i < max; i++ {
		if deltaX < 0 {
			deltaX++
		} else if deltaX > 0 {
			deltaX--
		}
		if deltaY < 0 {
			deltaY++
		} else if deltaY > 0 {
			deltaY--
		}

		regionId := models.CoordsToRegionId(x-deltaX, y-deltaY)
		region := world.GetRegion(regionId)

		baseX, baseY := region.GetBase()
		offsetX := x-deltaX - baseX
		offsetY := y-deltaY - baseY

		log.Printf("relx %d rely %d", offsetX, offsetY)
		tileKey := fmt.Sprintf("%d-%d-%d", offsetX, offsetY, 0)
		if _, ok := region.BlockedTiles[tileKey]; ok {
			log.Printf("tile is blocked")
			continue
		}

		e.addStep(x-deltaX, y-deltaY)
	}
}

func (e *movementEntity) addStep(x, y int) {
	last := e.getLast()
	deltaX := x - int(last.Tile.X)
	deltaY := y - int(last.Tile.Y)
	direction := models.DirectionFromDeltas(deltaX, deltaY)
	if direction != models.Direction.None {
		e.Steps = append(e.Steps, &models.Step{
			Tile: &models.Tile{
				X: uint16(x),
				Y: uint16(y),
			},
			Direction: direction,
		})
	}
}

func (e *movementEntity) getLast() *models.Step {
	var last *models.Step
	if len(e.Steps) > 0 {
		last = e.Steps[len(e.Steps)-1]
	} else {
		last = &models.Step{
			Tile:      e.Movement.Position,
			Direction: models.Direction.None,
		}
	}

	return last
}
