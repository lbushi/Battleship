// Battleship is a server that provides the functionality of the Battleship game.
// Battleship also provides the ability to have concurrent games which are independent of each other
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"unicode"
)

type Point struct {
	x, y int
}

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
		} else {
			return 'X'
		}
	} else {
		if hidden {
			return '.'
		} else {
			if cell.hasShip() {
				return cell.ship.name
			} else {
				return '.'
			}
		}
	}
}

func (cell *Cell) hasShip() bool {
	return cell.ship != nil
}

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

// the number of ships a player has(currently 1 for debugging purpose but will be around 4 or 5)
const shipNum int = 1

type Player struct {
	ships [shipNum]Ship // the list of ships that the player owns
}

// NewPlayer is a contructor function for the player type that initializes the ships for each player
func NewPlayer() *Player {
	var player Player
	var ship *Ship
	for i := 0; i < shipNum; i++ {
		ship = &player.ships[i]
		ship.setName(rune('A' + i))
		ship.appendPoint(Point{0, 0}, Point{0, 1})
	}
	return &player
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
//to initialize the players.
func NewGame() *Game {
	var game Game
	for i := 0; i < playerNum; i++ {
		game.players[i] = *NewPlayer()
	}
	for i := 0; i < playerNum; i++ {
		for j := 0; j < gridSize; j++ {
			for k := 0; k < gridSize; k++ {
				game.grids[i][j][k] = Cell{coordinates: Point{j, k}}
			}
		}
	}
	for i := 0; i < playerNum; i++ {
		for j, ship := range game.players[i].ships {
			for _, p := range ship.parts {
				game.grids[i][p.x][p.y].ship = &game.players[i].ships[j]
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

// Broadcast is a method on the minServer type that sends a message to all players of the server
func (minsrv *minServer) Broadcast(message string) {
	var wg sync.WaitGroup
	wg.Add(playerNum)
	for _, pl := range minsrv.players {
		go func(pl net.Conn) {
			defer wg.Done()
			io.WriteString(pl, message)
		}(pl)
	}
	wg.Wait()
}

// the main driver function for the long-lived communication between the server and the players
func (minsrv *minServer) Handle() {
	ch := make(chan struct{})
	minsrv.Broadcast("The game has started!\n")
	var wg sync.WaitGroup
	wg.Add(playerNum)
	for i, pl := range minsrv.players {
		go func(pl net.Conn) {
			defer wg.Done()
			io.WriteString(pl, minsrv.game.getDisplay(i))
		}(pl)
	}
	wg.Wait()
	for {
		go func() {
			io.WriteString(minsrv.players[minsrv.game.turn], "It is your turn!\n")
			ch <- struct{}{}
		}()
		go func() {
			io.WriteString(minsrv.players[1-minsrv.game.turn], "Your opponent is thinking...\n")
			ch <- struct{}{}
		}()
		// we are using ch to synchronize these two goroutines with the handle goroutine
		<-ch
		<-ch
		x, y := minsrv.GetInput(minsrv.game.turn)
		minsrv.game.Hit(x, y)
		wg.Add(2)
		for i, pl := range minsrv.players {
			go func(pl net.Conn, i int) {
				defer wg.Done()
				io.WriteString(pl, minsrv.game.getDisplay(i))
			}(pl, i)
		}
		wg.Wait()
		if minsrv.game.isOver() {
			minsrv.endGame(minsrv.game.turn)
			break
		}
		minsrv.game.changeTurn()
	}
}

// GetInput gets the input coordinates from the player whose turn it is at the moment
func (minsrv *minServer) GetInput(player int) (int, int) {
	var x, y int
	fmt.Fscan(minsrv.players[player], &x, &y)
	return x, y
}

// endGame sends the final win/loss message to the players and then closes the connections with each of them
func (minsrv *minServer) endGame(winner int) {
	for i, pl := range minsrv.players {
		go func(i int, pl net.Conn) {
			if i == winner {
				io.WriteString(pl, "You won!")
			} else {
				io.WriteString(pl, "You lost!")
			}
		}(i, pl)
	}
}

// Server struct represents the main server which will receive requests from players and will pair them with each other so that they can initiate the game
type Server struct {
	waitingcapacity int
	waitingConn     []net.Conn
}

func (srv *Server) Listen() error {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		if len(srv.waitingConn) == srv.waitingcapacity {
			srv.waitingConn = append(srv.waitingConn, conn)
			go (&minServer{game: NewGame(), players: srv.waitingConn}).Handle()
			srv.waitingConn = nil
		} else {
			fmt.Fprintf(conn, "You are player number %d!\nPlease wait for other players to join!\n", len(srv.waitingConn)+1)
			srv.waitingConn = append(srv.waitingConn, conn)
		}
	}
}

func main() {
	srv := new(Server)
	srv.waitingcapacity = playerNum - 1
	err := srv.Listen()
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: optimize the state calculation for a cell by adding a state field inside the cell and updating the state only when the cell is hit
