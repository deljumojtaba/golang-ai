package main

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
)

type DepthFirstSearch struct {
	Frontier []*Node
	Game     *Maze
}

func (dfs *DepthFirstSearch) GetFrontier() []*Node {
	return dfs.Frontier
}

func (dfs *DepthFirstSearch) AddToFrontier(node *Node) {
	dfs.Frontier = append(dfs.Frontier, node)
}

func (dfs *DepthFirstSearch) ContainsState(i *Node) bool {
	for _, node := range dfs.Frontier {
		if node.State == i.State {
			return true
		}
	}
	return false
}

func (dfs *DepthFirstSearch) IsEmpty() bool {
	return len(dfs.Frontier) == 0
}

func (dfs *DepthFirstSearch) RemoveFromFrontier() (*Node, error) {
	if !dfs.IsEmpty() {
		if dfs.Game.Debug {
			fmt.Printf("DFS removing node at %v from frontier\n", dfs.Frontier[len(dfs.Frontier)-1].State)
			for _, n := range dfs.Frontier {
				fmt.Printf(" - Node at %v\n", n.State)
			}
		}

		n := dfs.Frontier[len(dfs.Frontier)-1]
		dfs.Frontier = dfs.Frontier[:len(dfs.Frontier)-1]
		return n, nil
	}

	return nil, errors.New("frontier is empty")
}

func (dfs *DepthFirstSearch) Solve() {
	fmt.Println("Starting Depth-First Search")

	dfs.Game.NumExplored = 0

	start := &Node{
		State:  dfs.Game.Start,
		Parent: nil,
		Index:  0,
		Action: "",
	}

	dfs.AddToFrontier(start)
	dfs.Game.CurrentNode = start

	for {
		if dfs.IsEmpty() {
			fmt.Println("No solution found")
			return
		}

		currentNode, err := dfs.RemoveFromFrontier()
		if err != nil {
			fmt.Println("Error removing from frontier:", err)
			return
		}

		if dfs.Game.Debug {
			fmt.Printf("Exploring node at %v\n", currentNode.State)
			fmt.Println("----------------------------------------")
			fmt.Println()
		}

		dfs.Game.NumExplored++
		dfs.Game.CurrentNode = currentNode

		// Check if we have reached the goal
		if currentNode.State == dfs.Game.Goal {
			fmt.Println("Goal reached!")
			var actions []string
			var cells []Point

			for {
				if currentNode.Parent != nil {
					actions = append(actions, currentNode.Action)
					cells = append(cells, currentNode.State)
					currentNode = currentNode.Parent
				} else {
					break
				}
			}

			// Reverse the actions and cells to get the correct order
			slices.Reverse(actions)
			slices.Reverse(cells)

			dfs.Game.Solution = Solution{
				Actions: actions,
				Cells:   cells,
			}

			dfs.Game.Explored = append(dfs.Game.Explored, currentNode.State)

			// Stop search immediately after finding the goal.
			// Previously the code continued searching which could
			// print "Goal reached!" multiple times and later
			// print "No solution found" when the frontier emptied.
			break
		}

		dfs.Game.Explored = append(dfs.Game.Explored, currentNode.State)

		// Expand the current node to get its neighbors

		for _, neighbor := range dfs.Neighbors(currentNode) {
			if !dfs.ContainsState(neighbor) {
				if !inExplored(neighbor.State, dfs.Game.Explored) {
					dfs.AddToFrontier(&Node{
						State:  neighbor.State,
						Parent: currentNode,
						Action: neighbor.Action,
					})
				}
			}
		}
	}
}

func (dfs *DepthFirstSearch) Neighbors(node *Node) []*Node {
	row := node.State.Row
	col := node.State.Col

	candidates := []*Node{
		{State: Point{Row: row - 1, Col: col}, Parent: node, Action: "UP"},
		{State: Point{Row: row + 1, Col: col}, Parent: node, Action: "DOWN"},
		{State: Point{Row: row, Col: col - 1}, Parent: node, Action: "LEFT"},
		{State: Point{Row: row, Col: col + 1}, Parent: node, Action: "RIGHT"},
	}

	var neighbors []*Node
	for _, candidate := range candidates {
		if 0 <= candidate.State.Row && candidate.State.Row < dfs.Game.Height &&
			0 <= candidate.State.Col && candidate.State.Col < dfs.Game.Width &&
			!dfs.Game.Walls[candidate.State.Row][candidate.State.Col].wall {
			neighbors = append(neighbors, candidate)
		}
	}

	for i := range neighbors {
		j := rand.Intn(i + 1)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}

	return neighbors
}
