package main

type Player struct {
	ships [shipNum]Ship // the list of ships that the player owns
}

// addShip adds a ship with an orientation of dir(either vertical or horizontal) and coordinates of (x, y) to the set of ships of the player
func (player *Player) initShip(index, x, y int, dir string) {
	ship := &(player.ships[index])
	if dir == "horizontal" {
		ship.setName(rune('A' + index))
		for j := 0; j < shipSize; j++ {
			ship.appendPoint(Point{x, y + j})
		}
	} else {
		ship.setName(rune('A' + index))
		for j := 0; j < shipSize; j++ {
			ship.appendPoint(Point{x + j, y})
		}
	}

}

// isDefeated returns a boolean indicating whether all ships of a player have been destroyed
func (player *Player) isDefeated() bool {
	for _, ship := range player.ships {
		if !ship.IsDestroyed() {
			return false
		}
	}
	return true
}

// gridSize represents the board dimensions and playerNum represents the number of players who are playing
const (
	gridSize  int = 10
	playerNum int = 2
)
