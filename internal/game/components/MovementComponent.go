package components

import (
	"log"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/models"
	"rsps/util"
)

type MovementComponent struct {
	movement *models.Movement
	world    *game.World
	steps    []*models.Step
}

func NewMovementComponent(movement *models.Movement, world *game.World) *MovementComponent {
	return &MovementComponent{
		movement: movement,
		world:    world,
		steps:    make([]*models.Step, 0),
	}
}

func (m *MovementComponent) Tick() {
	if len(m.steps) == 0 {
		return
	}

	walkPoint := m.steps[0]
	m.steps = m.steps[1:]

	var runPoint *models.Step
	if m.movement.IsRunning && len(m.steps) > 0 {
		runPoint = m.steps[0]
		m.steps = m.steps[1:]
	}

	if walkPoint != nil {
		m.movement.Position = walkPoint.Tile
		m.movement.Direction = walkPoint.Direction
		log.Printf("wdir %+v", walkPoint.Direction)
	}

	if runPoint != nil {
		m.movement.Position = runPoint.Tile
		m.movement.Direction = runPoint.Direction
		log.Printf("rdir %+v", runPoint.Direction)
	}
}

func (m *MovementComponent) MoveTo(p *models.Tile) {
	m.steps = make([]*models.Step, 0)
	m.addPosition(p)
}

func (m *MovementComponent) addPosition(p *models.Tile) {
	last := m.getLast()
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

		m.addStep(x-deltaX, y-deltaY)
	}
}

func (m *MovementComponent) addStep(x, y int) {
	last := m.getLast()
	deltaX := x - int(last.Tile.X)
	deltaY := y - int(last.Tile.Y)
	direction := models.DirectionFromDeltas(deltaX, deltaY)
	if direction != models.Direction.None {

		// NOTE: Collision works based on our current tile + direction,
		chunk := m.world.GetOrLoadChunk(&models.Tile{X: last.Tile.X, Y: last.Tile.Y})
		// TODO: height instead of 0
		if chunk.CollisionMatrix[0].IsBlocked(int(last.Tile.X), int(last.Tile.Y), direction) {
			return
		}

		if direction.IsDiagonal() {
			for _, d := range direction.GetDiagonalComponents() {
				stepX, stepY := int(last.Tile.X) + d.GetDeltaX(), int(last.Tile.Y) + d.GetDeltaY()
				chunk := m.world.GetOrLoadChunk(&models.Tile{X: uint16(stepX), Y: uint16(stepY)})
				if chunk.CollisionMatrix[0].IsBlocked(stepX, stepY, d.GetOpposite()) {
					return
				}
			}
		}

		m.steps = append(m.steps, &models.Step{
			Tile: &models.Tile{
				X: uint16(x),
				Y: uint16(y),
			},
			Direction: direction,
		})
	}
}

func (m *MovementComponent) getLast() *models.Step {
	var last *models.Step
	if len(m.steps) > 0 {
		last = m.steps[len(m.steps)-1]
	} else {
		last = &models.Step{
			Tile:      m.movement.Position,
			Direction: models.Direction.None,
		}
	}

	return last
}
