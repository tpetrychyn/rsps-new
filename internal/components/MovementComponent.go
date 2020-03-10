package components

import (
	"log"
	"rsps-comm-test/pkg/models"
	"rsps/util"
)

type MovementComponent struct {
	Direction models.DirectionEnum
	Movement  *models.Movement
	Steps     []*models.Step
}

func NewMovementComponent(movement *models.Movement) *MovementComponent {
	return &MovementComponent{
		Movement: movement,
		Steps:    make([]*models.Step, 0),
	}
}

func (m *MovementComponent) Tick() {
	if len(m.Steps) == 0 {
		return
	}

	walkPoint := m.Steps[0]
	m.Steps = m.Steps[1:]

	var runPoint *models.Step
	if m.Movement.IsRunning && len(m.Steps) > 0 {
		runPoint = m.Steps[0]
		m.Steps = m.Steps[1:]
	}

	if walkPoint != nil  {
		m.Movement.Position = walkPoint.Position
		m.Movement.Direction = walkPoint.Direction
		log.Printf("wdir %+v", walkPoint.Direction)
	}

	if runPoint != nil {
		m.Movement.Position = runPoint.Position
		m.Movement.Direction = runPoint.Direction
		log.Printf("rdir %+v", runPoint.Direction)
	}
}

func (m *MovementComponent) MoveTo(p *models.Position) {
	m.Steps = make([]*models.Step, 0)
	m.addPosition(p)
}

func (m *MovementComponent) addPosition(p *models.Position) {
	last := m.getLast()
	x := int(p.X)
	z := int(p.Z)

	deltaX := x - int(last.Position.X)
	deltaZ := z - int(last.Position.Z)

	max := util.Abs(deltaX)
	if util.Abs(deltaZ) > util.Abs(deltaX) {
		max = util.Abs(deltaZ)
	}

	for i := 0; i < max; i++ {
		if deltaX < 0 {
			deltaX++
		} else if deltaX > 0 {
			deltaX--
		}
		if deltaZ < 0 {
			deltaZ++
		} else if deltaZ > 0 {
			deltaZ--
		}

		m.addStep(x-deltaX, z-deltaZ)
	}
}

func (m *MovementComponent) addStep(x, z int) {
	last := m.getLast()
	deltaX := x - int(last.Position.X)
	deltaY := z - int(last.Position.Z)
	direction := models.DirectionFromDeltas(deltaX, deltaY)
	if direction != models.Direction.None {
		m.Steps = append(m.Steps, &models.Step{
			Position: &models.Position{
				X: uint16(x),
				Z: uint16(z),
			},
			Direction: direction,
		})
	}
}

func (m *MovementComponent) getLast() *models.Step {
	var last *models.Step
	if len(m.Steps) > 0 {
		last = m.Steps[len(m.Steps)-1]
	} else {
		last = &models.Step{
			Position: m.Movement.Position,
			Direction: models.Direction.None,
		}
	}

	return last
}
