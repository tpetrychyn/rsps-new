package components

import (
	"github.com/tpetrychyn/rsps-comm-test/internal/game"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"log"
)

type MovementComponent struct {
	movement        *models.Movement
	world           *game.World
	steps           []*models.Step
	addMovementFlag func()
	setMapFlag      func(*models.Tile)
}

func NewMovementComponent(movement *models.Movement, world *game.World, addMovementFlag func(), setMapFlag func(*models.Tile)) *MovementComponent {
	return &MovementComponent{
		movement:        movement,
		world:           world,
		steps:           make([]*models.Step, 0),
		addMovementFlag: addMovementFlag,
		setMapFlag:      setMapFlag,
	}
}

func (m *MovementComponent) Clear() {
	m.steps = make([]*models.Step, 0)
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
		if m.IsBlocked(walkPoint, runPoint) || m.IsBlocked(runPoint, walkPoint) {
			log.Printf("tried to run on blocked step")
		}
		m.steps = m.steps[1:]
	}

	if walkPoint != nil {
		m.movement.Position = walkPoint.Tile
		m.movement.WalkDirection = walkPoint.Direction
	}

	if runPoint != nil {
		m.movement.Position = runPoint.Tile
		m.movement.RunDirection = runPoint.Direction
		m.addMovementFlag()
	}
}

func (m *MovementComponent) MoveTo(p *models.Tile) {
	m.steps = make([]*models.Step, 0)
	m.calculateRoute(p)
}

func (m *MovementComponent) calculateRoute(dest *models.Tile) {
	nodes := make([]*models.Step, 0)
	visited := make([]*models.Step, 0)

	start := m.getLast()
	if start.X == dest.X && start.Y == dest.Y {
		return
	}

	nodes = append(nodes, start)
	for {
		if len(nodes) == 0 {
			break
		}
		head := nodes[0]
		nodes = nodes[1:]

		for _, dir := range models.RsDirectionOrder {
			tile := &models.Step{Tile: head.Step(dir), Direction: dir, Head: head, Cost: head.Cost + 1}
			if !start.IsWithinRadius(tile.Tile, 20) {
				continue
			}
			if containsTile(visited, tile) {
				continue
			}

			if m.IsBlocked(head, tile) || m.IsBlocked(tile, head) {
				continue
			}
			nodes = append(nodes, tile)
			visited = append(visited, tile)
			if tile.X == dest.X && tile.Y == dest.Y {
				m.setMapFlag(dest)
				for {
					if tile.X == start.X && tile.Y == start.Y {
						return
					}
					deltaX := tile.X - tile.Head.X
					deltaY := tile.Y - tile.Head.Y
					direction := models.DirectionFromDeltas(deltaX, deltaY)
					tile.Direction = direction
					m.steps = append([]*models.Step{tile}, m.steps...)
					tile = tile.Head
				}
			}
		}
	}
	// couldnt make it to tile, see if we can get closer
	min := start
	for _, v := range visited {
		if v.Tile.DistanceTo(dest) < min.Tile.DistanceTo(dest) {
			min = v
		}
	}
	m.calculateRoute(min.Tile)
}

func containsTile(arr []*models.Step, tile *models.Step) bool {
	for _, v := range arr {
		if v.X == tile.X && v.Y == tile.Y {
			return true
		}
	}
	return false
}

func (m *MovementComponent) IsBlocked(from, dest *models.Step) bool {
	deltaX := dest.X - from.X
	deltaY := dest.Y - from.Y
	direction := models.DirectionFromDeltas(deltaX, deltaY)
	chunk := m.world.GetOrLoadChunk(from.Tile)
	if chunk.CollisionMatrix[m.movement.Position.Height].IsBlocked(from.X, from.Y, direction) {
		return true
	}

	if direction.IsDiagonal() {
		for _, d := range direction.GetDiagonalComponents() {
			step := from.Step(d)
			chunk := m.world.GetOrLoadChunk(step)
			if chunk.CollisionMatrix[m.movement.Position.Height].IsBlocked(step.X, step.Y, d.GetOpposite()) {
				return true
			}
		}
	}
	return false
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
