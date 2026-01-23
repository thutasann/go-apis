package elevator

type Direction int

const (
	Idle Direction = iota
	Up
	Down
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "UP"
	case Down:
		return "DOWN"
	default:
		return "IDLE"
	}
}
