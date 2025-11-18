package main

import (
	"fmt"
	"os"
)

func inExplored(needle Point, haystack []Point) bool {
	for _, p := range haystack {
		if p.Col == needle.Col && p.Row == needle.Row {
			return true
		}
	}
	return false
}

func emptyTmp() {
	files, err := os.ReadDir("./tmp")
	if err != nil {
		fmt.Printf("error reading tmp directory: %v", err)
		return
	}

	for _, file := range files {
		err := os.Remove("./tmp/" + file.Name())
		if err != nil {
			fmt.Printf("error removing file %s: %v", file.Name(), err)
		}
	}
}
