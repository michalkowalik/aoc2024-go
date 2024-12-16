package main

type Position struct {
	x, y int
}

type Guardian struct {
	position  Position
	direction int
	path      []Position // save the path the guardian follows
	inTheLoop bool
}

func NewGuardian(x int, y int, direction int) *Guardian {
	guardian := new(Guardian)
	guardian.inTheLoop = false
	guardian.position.x = x
	guardian.position.y = y
	guardian.direction = direction
	guardian.path = append(guardian.path, guardian.position)
	return guardian
}

func (w *Guardian) Move(storage [][]MapObject) {
	nextPosition := Position{w.position.x, w.position.y}
	// next position:
	switch w.direction {
	case 0:
		nextPosition.y--
	case 1:
		nextPosition.x++
	case 2:
		nextPosition.y++
	default:
		nextPosition.x--
	}

	if _, ok := (storage[nextPosition.y][nextPosition.x]).(*Obstacle); ok {
		// we've bounced from the obstacle:
		inTheLoop := storage[nextPosition.y][nextPosition.x].Bounce(w.direction)

		if inTheLoop {
			w.inTheLoop = true
			return
		}

		// and we change the direction:
		w.direction = (w.direction + 1) % 4
	} else {
		w.position = nextPosition
		w.path = append(w.path, w.position)
	}
}

func (w *Guardian) ChangeDirection() {

}

func (w *Guardian) Path() []Position {
	return w.path
}

func (w *Guardian) isLeavingStorage() bool {
	if w.position.x == 0 && w.direction == 3 {
		return true
	}
	if w.position.y == 0 && w.direction == 0 {
		return true
	}
	if w.position.x == storageSize-1 && w.direction == 1 {
		return true
	}
	if w.position.y == storageSize-1 && w.direction == 2 {
		return true
	}
	return false
}
