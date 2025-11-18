package main

// maze-ai is a small program that loads a text-based maze from a file
// and (eventually) will run different graph search algorithms on it
// (DFS, BFS, GBFS, A*, Dijkstra). This file currently focuses on
// parsing and representing the maze in memory.

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Search algorithm identifiers. These constants are used to select
// which search algorithm to run when the solver is implemented.
const (
	DFS = iota
	BFS
	GBFS
	ASTAR
	DIJKSTRA
)

// Point is a simple row/column coordinate in the maze grid.
// Row and Col are zero-indexed and correspond to the line and
// character position from the input file.
type Point struct {
	Row int
	Col int
}

// Wall represents a single cell in the maze. The field name `wall`
// is true when the cell is a blocking wall (#). When false the cell
// is walkable (space, start 'A' or goal 'B'). State stores the cell's
// coordinates; this is redundant with the 2D slice indexes but kept
// for convenience.
type Wall struct {
	State Point
	wall  bool
}

type Node struct {
	State  Point
	Parent *Node
	Index  int
	Action string
}

type Solution struct {
	Actions []string
	Cells   []Point
}

// Maze stores the parsed maze. Height and Width reflect the number
// of lines and the number of characters per line (note: Width
// currently uses the length of the first line and may include the
// newline character if not trimmed). Start and Goal are points for
// 'A' and 'B', respectively. Walls is a 2D grid of Wall cells.
type Maze struct {
	Height      int
	Width       int
	Start       Point
	Goal        Point
	Walls       [][]Wall
	CurrentNode *Node
	Solution    Solution
	Explored    []Point
	Steps       int
	NumExplored int
	Debug       bool
	SearchType  int
	Animate     bool
}

func init() {
	// Ensure the tmp directory is empty at the start of the program.
	_ = os.Mkdir("./tmp", os.ModePerm)
	emptyTmp()
}

func main() {
	// Example usage
	var m Maze
	var maze, searchType string

	flag.StringVar(&maze, "file", "maze.txt", "maze file")
	flag.StringVar(&searchType, "search", "BFS", "search type")
	flag.BoolVar(&m.Debug, "debug", false, "enable debug mode")
	flag.BoolVar(&m.Animate, "animate", false, "generate frames (PNG) in ./tmp instead of a final animation")
	flag.Parse()

	// Print the chosen file and search type (for now the program
	// only loads the maze; search algorithms are to be wired later).
	fmt.Println("Maze file:", maze)
	fmt.Println("Search type:", searchType)

	// Load the maze from the file. Load fills the Maze struct or
	// returns an error if the file can't be read or required
	// characters (A/B) are missing.
	err := m.Load(maze)
	if err != nil {
		fmt.Println("Error loading maze:", err)
		os.Exit(1)
	}

	startTime := time.Now()

	switch searchType {
	case "dfs":
		m.SearchType = DFS
		solveDFS(&m)
	case "bfs":
		m.SearchType = BFS
	case "gbfs":
		m.SearchType = GBFS
	case "astar":
		m.SearchType = ASTAR
	case "dijkstra":
		m.SearchType = DIJKSTRA
	default:
		fmt.Println("Unknown search type:", searchType)
		os.Exit(1)
	}

	if len(m.Solution.Actions) > 0 {
		fmt.Println("Solution found:")
		// m.printMaze()
		fmt.Println("Steps to goal:", len(m.Solution.Cells))
		fmt.Println("Time taken:", time.Since(startTime))
		m.OutputImage("image.png")
	} else {
		fmt.Println("No solution found.")
	}

	fmt.Println("Total nodes explored:", m.NumExplored)

	if m.Animate {
		// create final animation from frames in ./tmp
		m.OutputAnimatedImage()
		fmt.Printf("Frames written to ./tmp (%d frames).\n", m.NumExplored)
	}
}

func solveDFS(m *Maze) {
	var dfs DepthFirstSearch
	dfs.Game = m
	fmt.Println("Goal is at:", dfs.Game.Goal)
	dfs.Solve()
}

func (g *Maze) printMaze() {
	for r, row := range g.Walls {
		for c, col := range row {
			if col.wall {
				fmt.Print("#")
			} else if g.Start.Row == col.State.Row && g.Start.Col == col.State.Col {
				fmt.Print("A")
			} else if g.Goal.Row == col.State.Row && g.Goal.Col == col.State.Col {
				fmt.Print("B")
			} else if g.inSolution(Point{r, c}) {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Maze) inSolution(x Point) bool {
	for _, step := range g.Solution.Cells {
		if step.Col == x.Col && step.Row == x.Row {
			return true
		}
	}
	return false
}

func (g *Maze) Load(filename string) error {
	// Open the file for reading. Any error is returned to the caller.
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening %s: %v", filename, err)
		return err
	}
	defer f.Close()

	// Read the file line-by-line into a slice. Note: ReadString('\n')
	// keeps the trailing newline, so lines will often end with '\n'.
	// Depending on how you want Width to behave you may want to
	// strings.TrimRight(line, "\r\n") each line.
	var fileContents []string

	// Use a Scanner to read lines without trailing newlines. This
	// ensures Width corresponds to the number of visible characters
	// per line and prevents mismatches between Width and the
	// number of parsed cells (which caused index out-of-range
	// panics when printing the maze).
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fileContents = append(fileContents, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading %s: %v", filename, err)
		return err
	}

	// Quick check: ensure both start (A) and goal (B) exist in the
	// file. We stop scanning early if both are found.
	foundStart, foundEnd := false, false
	for _, line := range fileContents {
		if strings.Contains(line, "A") {
			foundStart = true
		}
		if strings.Contains(line, "B") {
			foundEnd = true
		}
		if foundStart && foundEnd {
			break
		}
	}

	if !foundStart || !foundEnd {
		return fmt.Errorf("maze must have start (A) and end (B) points")
	}

	// Set Height/Width. Using Scanner above yields lines without
	// trailing newlines so Width will match the number of
	// characters parsed into each row.
	g.Height = len(fileContents)
	if g.Height > 0 {
		g.Width = len(fileContents[0])
	} else {
		g.Width = 0
	}

	var rows [][]Wall

	// Parse each character in each line and build a grid of Walls.
	// Note: iterating over a string yields runes; using the index
	// 'j' with the rune value matches byte positions for ASCII
	// mazes (which is typical here).
	for i, row := range fileContents {
		var cols []Wall

		for j, col := range row {
			curLetter := fmt.Sprintf("%c", col)

			var wall Wall

			switch curLetter {
			case "A":
				// Start position
				g.Start = Point{i, j}
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case "B":
				// Goal position
				g.Goal = Point{i, j}
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case " ":
				// Walkable space
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = false
			case "#":
				// Wall / obstacle
				wall.State.Row = i
				wall.State.Col = j
				wall.wall = true
			default:
				// Ignore any other characters (for example a
				// trailing newline will produce '\n' here and is
				// skipped). This also prevents non-maze characters
				// from creating entries.
				continue
			}

			cols = append(cols, wall)

		}

		rows = append(rows, cols)

	}

	g.Walls = rows
	return nil

}
