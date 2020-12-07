# Battleship
Multiplayer Battleship Game implemented in Golang

Battleship is a concurrent multi-player command line implementation of the original Battleship game with board dimensions of 10x10 and where each player has 5 ships with the slight change that each ship occupies 5 cells on the board. The game has networking capabilities so that it allows two clients/players in different computers but in the same LAN to play with each other. 

Simply clone the repository on your own machine, run go build  and then start the produced executable which will wait for incoming requests from players who want to play and will pair them with other players. To request a game as a player you can use the nc command in UNIX-like systems to connect to the server on port 8000, given that the server is already running on the target computer. 

Specify the ship positions by providing for each ship the leftmost/topmost cell and then providing a leter equal to h/v for whether the ship is oriented horizontally or vertically. To bomb a cell, provide the cell coordinates ranging from 0-9 and separated by space when its your turn.
