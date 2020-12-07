package main

import (
	"net"
	"time"
)

// the Game struct represent the model of the game
type Game struct {
	turn    int
	grids   [playerNum][gridSize][gridSize]Cell
	players [playerNum]Player
}

// getPlayerDisplay returns the partial display board of a player together with a parameter which determines whether the cells will be revealed or not
func (game *Game) getPlayerDisplay(playerid int, hidden bool) (display string) {
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			display += string(game.grids[playerid][i][j].getState(hidden))
			if j != gridSize-1 {
				display += " "
			}
		}
		display += "\n"
	}
	return display
}

func (game *Game) changeTurn() {
	game.turn = (game.turn + 1) % playerNum
}

func (game *Game) Hit(x, y int) {
	game.grids[1-game.turn][x][y].hit()
}

// getDisplay returns the total display board of a player, that is, it includes both the player's and the player's opponents displays
func (game *Game) getDisplay(playerid int) (display string) {
	var midline string
	for i := 0; i < 2*gridSize; i++ {
		midline += "-"
	}
	midline += "\n"
	return game.getPlayerDisplay(playerid, false) + midline + game.getPlayerDisplay(1-playerid, true)
}

// isOver returns a boolean indicating whether the game is over
func (game *Game) isOver() bool {
	return game.players[1-game.turn].isDefeated()
}

// NewGame is a constructor function that initializes the game by setting up the relations between the players, cells, and ships. It calls NewPlayer, defined above,
// to initialize the players.
func NewGame() *Game {
	var game Game
	for i := 0; i < playerNum; i++ {
		game.players[i] = Player{}
	}
	for i := 0; i < playerNum; i++ {
		for j := 0; j < gridSize; j++ {
			for k := 0; k < gridSize; k++ {
				game.grids[i][j][k] = Cell{coordinates: Point{j, k}}
			}
		}
	}
	return &game
}

type minServer struct {
	game    *Game
	timeout time.Duration // the duration that the server will wait for a read from a connection before it moves on(currently not used but will add later)
	players []net.Conn
}
