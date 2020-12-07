package main

import (
	"unicode"
)

// Point is a struct representing a pair of coordinates on the plane
type Point struct {
	x, y int
}

// Cell represents a cell on the battleship board
type Cell struct {
	coordinates Point // the cell coordinates
	ship        *Ship // the ship that has a piece on this cell, if it exists, otherwise it is nil
	Hit         bool  // boolean that indicates whether the cell has been hit
}

// this functions hits the cell
func (cell *Cell) hit() {
	if cell.ship != nil {
		cell.ship.hit()
	}
	cell.Hit = true
}

// this function return true if and only if the cell has been hit
func (cell *Cell) isHit() bool {
	return cell.Hit
}

// getState returns the character that will represent the cell in the display based on whether the cell has been hit and whether we want to hide it
func (cell *Cell) getState(hidden bool) rune {
	if cell.isHit() {
		if cell.hasShip() {
			return unicode.ToLower(cell.ship.name)
		}
		return 'X'
	}
	if hidden {
		return '.'
	}
	if cell.hasShip() {
		return cell.ship.name
	}
	return '.'

}

func (cell *Cell) hasShip() bool {
	return cell.ship != nil
}
