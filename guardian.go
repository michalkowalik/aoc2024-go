package main

type Position struct {
	x, y int
}

var moveCount int

type Guardian struct {
	position  Position
	direction int
	path      []Position // save the path the guardian follows
}

func NewGuardian(x int, y int, direction int) *Guardian {
	guardian := new(Guardian)
	guardian.position.x = x
	guardian.position.y = y
	guardian.direction = direction
	guardian.path = append(guardian.path, guardian.position)
	return guardian
}

func (w *Guardian) Move(storage StorageMap) {
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

	if _, ok := (storage.MapArray[nextPosition.y][nextPosition.x]).(*Obstacle); ok {
		w.direction = (w.direction + 1) % 4
	} else {
		w.position = nextPosition
		w.path = append(w.path, w.position)
		moveCount++
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
