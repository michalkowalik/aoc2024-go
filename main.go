package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type StorageMap struct {
	MapArray [][]MapObject
}

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
	storage.MapArray = make([][]MapObject, 0)

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
		storage.MapArray = append(storage.MapArray, row)
		lineIndex++
	}
	storageSize = len(storage.MapArray)
	run(guardian, storage)
}

func run(guardian *Guardian, storage StorageMap) {
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

func printCompletedPath(storageMap StorageMap, guardian *Guardian) {
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
			if _, ok := (storageMap.MapArray[y][x]).(*Obstacle); ok {
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
