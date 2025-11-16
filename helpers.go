package main

func inExplored(needle Point, haystack []Point) bool {
	for _, p := range haystack {
		if p.Col == needle.Col && p.Row == needle.Row {
			return true
		}
	}
	return false
}
