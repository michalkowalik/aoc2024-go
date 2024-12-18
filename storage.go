package main

import "sync"

type MapObject interface {
	Bounce(direction int) bool
}

type EmptySquare struct{}

type Obstacle struct {
	// every time the guardian bounces the obstacle, save the direction they came from.
	bounces map[int]bool
}

type StorageMap struct {
	sync.RWMutex
	storageMap [][]MapObject
}

func (e EmptySquare) Bounce(_ int) bool {
	return false
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
	if _, v := o.bounces[direction]; !v {
		o.bounces[direction] = true
		return false
	} else {
		return true
	}
}

func (s *StorageMap) CopyStorageMap() [][]MapObject {
	s.RLock()
	defer s.RUnlock()
	// Create a new slice of slices with the same dimensions
	newStorage := make([][]MapObject, len(s.storageMap))

	for i := range s.storageMap {
		newStorage[i] = append([]MapObject{}, s.storageMap[i]...)

		// remove any information about bounces
		for j, mapObject := range newStorage[i] {
			if _, ok := mapObject.(*Obstacle); ok {
				newStorage[i][j] = NewObstacle()
			}
		}
	}

	return newStorage
}
