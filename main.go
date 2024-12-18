package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

var storageSize = 0

//goland:noinspection GoDfaNilDereference
func main() {
	var guardian *Guardian
	f, err := os.Open("input-day6.txt")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	r := bufio.NewReader(f)
	storage := StorageMap{}
	storage.storageMap = make([][]MapObject, 0)

	lineIndex := 0
	for {
		line, _, err := r.ReadLine()

		if err == io.EOF {
			break
		}
		columnIndex := 0
		row := make([]MapObject, 0)
		for _, x := range line {
			if string(x) == "#" {
				row = append(row, NewObstacle())
			} else if string(x) == "." {
				row = append(row, EmptySquare{})
			} else if string(x) == "^" {
				guardian = NewGuardian(columnIndex, lineIndex, 0)
				row = append(row, EmptySquare{})
			}
			columnIndex++
		}
		storage.storageMap = append(storage.storageMap, row)
		lineIndex++
	}
	storageSize = len(storage.storageMap)
	run(guardian, storage.storageMap)

	// complete path will be needed for the part 2
	part2(*guardian, storage)
}

func run(guardian *Guardian, storage [][]MapObject) {
	fmt.Printf("Guardian starts at %d, %d, direction %d\n",
		guardian.position.x, guardian.position.y, guardian.direction)

	// start moving around....
	for !guardian.isLeavingStorage() {
		guardian.Move(storage)
	}

	fmt.Printf("Guardian's path is %d steps long\n",
		len(guardian.Path()))

	guardian.PrintCompletedPath(storage)
}

// part 2
func part2(guardian Guardian, storage StorageMap) {
	startingPosition := guardian.Path()[0]
	loopCounter := 0
	path := guardian.Path()[1:]
	ch := make(chan bool)

	for _, pos := range path {
		go runGuardianWithObstacle(storage, pos, startingPosition, ch)
	}

	for range path {
		result := <-ch
		if result {
			loopCounter++
		}
	}
	fmt.Printf("\nPART2: Number of loops: %d\n", loopCounter)
}

// return true if adding obstacle caused the loop
func runGuardianWithObstacle(storage StorageMap, obstacle Position, startingPosition Position, ch chan bool) {
	copiedStorage := storage.CopyStorageMap()

	guardian := NewGuardian(startingPosition.x, startingPosition.y, 0)
	copiedStorage[obstacle.y][obstacle.x] = NewObstacle()

	for !guardian.isLeavingStorage() {
		guardian.Move(copiedStorage)
		if guardian.inTheLoop {

			// why is that needed - where is the problem?
			time.Sleep(500 * time.Millisecond)
			ch <- true
		}
	}
	ch <- false
}
