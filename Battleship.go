// Battleship is a server that provides the functionality of the Battleship game.
// Battleship also provides the ability to have concurrent games which are independent of each other
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

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
	var wg sync.WaitGroup
	wg.Add(playerNum)
	for i, pl := range minsrv.players {
		go func(i int, pl net.Conn) {
			defer wg.Done()
			io.WriteString(pl, "Please provide the ship positions for each of your 5 ships!\n")
			var x, y int
			var dir string
			for j := 0; j < shipNum; j++ {
				fmt.Fscan(pl, &x, &y, &dir)
				minsrv.game.players[i].initShip(j, x, y, dir)
			}
			for i := 0; i < playerNum; i++ {
				for j, ship := range minsrv.game.players[i].ships {
					for _, p := range ship.parts {
						minsrv.game.grids[i][p.x][p.y].ship = &minsrv.game.players[i].ships[j]
					}
				}
			}
		}(i, pl)
	}
	wg.Wait()
	minsrv.Broadcast("The game has started!\n")
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
		fmt.Fprintf(minsrv.players[1-minsrv.game.turn], "Your opponent hit cell %d, %d!\n", x, y)
		if minsrv.game.isOver() {
			minsrv.endGame(minsrv.game.turn)
			break
		}
		minsrv.game.changeTurn()
	}
}

// GetInput gets the input coordinates from the player whose turn it is at the moment or it gets the ship positions if the game has just started and the players
// must initialize their ships
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
