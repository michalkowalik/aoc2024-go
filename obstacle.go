package main

type MapObject interface {
	Bounce(direction int) bool
}

type EmptySquare struct{}

func (e EmptySquare) Bounce(_ int) bool {
	return false
}

type Obstacle struct {
	// every time the guardian bounces the obstacle, save the direction they came from.
	bounces map[int]bool
}

func NewObstacle() *Obstacle {
	obstacle := &Obstacle{}
	obstacle.bounces = make(map[int]bool)
	return obstacle
}

// Bounce : Guardian bounces the obstacle.
// if the obstacle has been bounced already from the same direction, it means
// we're in the loop.
// return true if loop detected, false otherwise.
func (o *Obstacle) Bounce(direction int) bool {
	if !setBounces {
		return false
	}
	if _, v := o.bounces[direction]; !v {
		o.bounces[direction] = true
		return false
	} else {
		return true
	}
}
