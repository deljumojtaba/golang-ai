# maze-ai

A small Go program that loads a text-based maze from a file and can solve it using search algorithms (DFS implemented; BFS/GBFS/A*/Dijkstra are planned).

## Project structure

- `main.go` — maze loader and CLI
- `dfs.go` — depth-first search implementation
- `maze.txt`, `maze2.txt`, `maze-100-steps.txt` — example maze files

## Requirements

- Go 1.18+ (module-aware; tested with recent Go versions)

## Build

From the project root run:

```bash
go build
```

This produces a binary in the current directory (name depends on your OS).

## Run

You can run the program directly with `go run` or run the built binary.

Example (DFS):

```bash
go run . -file=maze.txt -search=dfs
```

Flags:
- `-file` : path to the maze file (default: `maze.txt`)
- `-search` : search algorithm to use. Supported values (case-insensitive): `dfs`, `bfs`, `gbfs`, `astar`, `dijkstra`. Only `dfs` is implemented by default.

Maze format:
- Characters per cell:
  - `A` — start
  - `B` — goal
  - `#` — wall/obstacle
  - space (` `) — walkable cell
- Each line represents a row of the maze. Keep lines the same length.

## Example output

```
Maze file: maze.txt
Search type: dfs
Goal is at: {0 2}
Starting Depth-First Search
Solution found:
... (maze printed with '.' for solution path)
Steps to goal: 10
Time taken: 1.234ms
Total nodes explored: 42
```

## Notes & suggestions

- The loader currently expects lines to be rectangular; if your maze file is ragged, the loader may error or printing may panic. The loader uses trimmed lines so trailing newlines do not cause width mismatches.
- To add more algorithms, follow the pattern used in `DepthFirstSearch` and wire them in `main.go`'s switch on `-search`.
- Feel free to open an issue or request if you'd like me to implement BFS or visualize the path.

## License

This repo contains example code; add a license file if you intend to share publicly.
