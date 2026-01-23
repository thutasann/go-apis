package elevator

type Elevator struct {
	currentFloor int
	direction    Direction
	targets      []int
}

func New(startFloor int) *Elevator {
	return &Elevator{
		currentFloor: startFloor,
		direction:    Idle,
		targets:      []int{},
	}
}

// adds a floot to the target queue
func (e *Elevator) AddRequest(floot int) {
	e.targets = append(e.targets, floot)

	if e.direction == Idle {
		e.updateDirection()
	}
}

// Step simulates one unit of time
func (e *Elevator) Step() {
	if len(e.targets) == 0 {
		e.direction = Idle
		return
	}

	target := e.targets[0]

	switch {
	case e.currentFloor < target:
		e.direction = Up
		e.currentFloor++
	case e.currentFloor > target:
		e.direction = Down
		e.currentFloor--
	default:
		e.targets = e.targets[1:]
		e.updateDirection()
	}
}

func (e *Elevator) updateDirection() {
	if len(e.targets) == 0 {
		e.direction = Idle
		return
	}

	target := e.targets[0]

	if target > e.currentFloor {
		e.direction = Up
	} else if target < e.currentFloor {
		e.direction = Down
	} else {
		e.direction = Idle
	}
}

func (e *Elevator) CurrentFloor() int {
	return e.currentFloor
}

func (e *Elevator) Direction() Direction {
	return e.direction
}

func (e *Elevator) Targets() []int {
	return append([]int(nil), e.targets...)
}

func (e *Elevator) IsIdle() bool {
	return e.direction == Idle && len(e.targets) == 0
}
