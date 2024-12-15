package main

type MapObject interface {
	Bounce(direction int)
	Visitations() []int
}

type EmptySquare struct{}

func (e EmptySquare) Bounce(_ int) {
	// empty implementation
}

func (e EmptySquare) Visitations() []int {
	return nil
}

type Obstacle struct {
	bounces []int // every time the guardian bounces the obstacle, save the direction they came from
}

func NewObstacle() *Obstacle {
	return &Obstacle{}
}

func (o *Obstacle) Bounce(direction int) {
	o.bounces = append(o.bounces, direction)
}

func (o *Obstacle) Visitations() []int {
	return o.bounces
}
