package main

// the number of ships a player has(currently 1 for debugging purpose but will be around 4 or 5)
const (
	shipNum  int = 5
	shipSize int = 4
)

type Ship struct {
	name     rune // the name of the ship which can be either A, B, C or D
	parts    []Point
	timesHit int // the number of cells of the ship which have been hit
}

func (ship *Ship) appendPoint(p ...Point) {
	ship.parts = append(ship.parts, p...)
}

func (ship *Ship) getName() rune {
	return ship.name
}

func (ship *Ship) IsDestroyed() bool {
	return ship.timesHit == len(ship.parts)
}

func (ship *Ship) hit() {
	ship.timesHit++
}

func (ship *Ship) setName(name rune) {
	ship.name = name
}
