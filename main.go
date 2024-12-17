package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var storage [][]MapObject
var setBounces = false
var storageSize = 0

//goland:noinspection GoDfaNilDereference
func main() {
	var guardian *Guardian
	f, err := os.Open("input-test.txt")
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

	storage = make([][]MapObject, 0)

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
		storage = append(storage, row)
		lineIndex++
	}
	storageSize = len(storage)
	run(guardian, storage)

	// complete path will be needed for the part 2
	part2(*guardian)
}

func run(guardian *Guardian, storage [][]MapObject) {
	fmt.Printf("Guardian starts at %d, %d, direction %d\n",
		guardian.position.x, guardian.position.y, guardian.direction)

	// start moving around....
	for !guardian.isLeavingStorage() {
		guardian.Move(storage)
	}

	fmt.Printf("Guardian's path is %d steps long\n",
		len(removeDuplicatesFromPath(guardian.Path())))

	printCompletedPath(storage, guardian)
}

// printCompletedPath visualizes the completed path of the guardian on the storage map including obstacles and prints it.
func printCompletedPath(storageMap [][]MapObject, guardian *Guardian) {
	stringMap := make([][]string, storageSize)

	for y := 0; y < storageSize; y++ {
		row := make([]string, storageSize)
		for x := 0; x < storageSize; x++ {
			row[x] = "."
		}
		stringMap[y] = row
	}

	for _, position := range guardian.Path() {
		stringMap[position.y][position.x] = "X"
	}

	for y := 0; y < storageSize; y++ {
		for x := 0; x < storageSize; x++ {
			if _, ok := (storageMap[y][x]).(*Obstacle); ok {
				stringMap[y][x] = "#"
			}
		}
	}

	for _, row := range stringMap {
		fmt.Printf("%s\n", strings.Join(row, ""))
	}
}

func removeDuplicatesFromPath(path []Position) []Position {
	allKeys := make(map[Position]bool)
	dedupPath := make([]Position, 0)
	for _, pos := range path {
		if _, v := allKeys[pos]; !v {
			allKeys[pos] = true
			dedupPath = append(dedupPath, pos)
		}
	}
	return dedupPath
}

// part 2
func part2(guardian Guardian) {
	startingPosition := guardian.Path()[0]

	loopCounter := 0

	dedupedPath := removeDuplicatesFromPath(guardian.Path())[1:]
	ch := make(chan bool, len(dedupedPath))

	for i := 1; i < len(dedupedPath); i++ {
		go runGuardianWithObstacle(dedupedPath[i], startingPosition, ch)
	}

	for i := 0; i < len(dedupedPath); i++ {
		result := <-ch
		if result {
			loopCounter++
		}
	}

	fmt.Printf("\nPART2: Number of loops: %d\n", loopCounter)
}

func copyMapObjects(original [][]MapObject) [][]MapObject {
	fmt.Printf("copying map objects\n")
	copied := make([][]MapObject, len(original))

	for i := range original {
		copied[i] = make([]MapObject, len(original[i]))
		copy(copied[i], original[i]) // Use copy for inner slices

		// remove any information about bounces
		for j, obstacle := range copied[i] {
			if _, ok := obstacle.(*Obstacle); ok {
				copied[i][j] = NewObstacle()
			}
		}
	}
	return copied
}

// return true if adding obstacle caused the loop
func runGuardianWithObstacle(obstacle Position, startingPosition Position, ch chan bool) {
	var storageLock sync.Mutex
	storageLock.Lock()
	copiedStorage := copyMapObjects(storage)
	storageLock.Unlock()

	setBounces = true
	guardian := NewGuardian(startingPosition.x, startingPosition.y, 0)
	copiedStorage[obstacle.y][obstacle.x] = NewObstacle()

	for !guardian.isLeavingStorage() {
		guardian.Move(copiedStorage)
		if guardian.inTheLoop {
			fmt.Printf("loop detected at %d, %d\n", obstacle.x, obstacle.y)
			printCompletedPath(copiedStorage, guardian)
			ch <- true
		}
	}
	ch <- false
}
