package main

/*
After writing this file with very little thought, the code in here has become very messy and gross
In my defense, everything's kinda magic number-y in how it needs to be done, so a programmatic approach was not obvious
*/

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var rowlens = [...]int{
	7,
	9,
	11,
	11,
	9,
	7,
}
var rowoff = [...]int{
	0,
	3,
	7,
	7,
	12,
	16,
}

var diceOdds = [...]float32{
	0,
	.028,
	.056,
	.083,
	.111,
	.139,
	.167,
	.139,
	.111,
	.083,
	.056,
	.028,
}

type Tile struct {
	roll   int
	bandit bool
	res    string
}

type Node struct {
	index int

	building int
	owner    int
	edges    []*Edge
	tiles    []*Tile
}

type Edge struct {
	inUse bool

	index int
	nodes []*Node
	road  int
}

type Board struct {
	tiles []*Tile
	nodes []*Node
	edges []*Edge
}

type Player struct {
	human         bool
	victoryPoints int
	number        int

	villages []*Node
	roads    []*Edge

	res   map[string]int
	cards []string

	knightsPlayed int
}

type Move struct {
	newNodes []Node
	newEdges []Edge

	cost        map[string]int
	card        string
	cardsBought int

	heur float32
}

type Game struct {
	players [4]*Player
	board   Board

	cards [26]string

	longestRoad int
	lrPlayer    int
	largestArmy int
	laPlayer    int
}

// We don't talk about this function
// Just act like it magically works
func NewBoard() Board {
	m := Board{}
	tiles := [...]string{
		"L", // Lumber
		"L",
		"L",
		"L",
		"W", // Wheat
		"W",
		"W",
		"W",
		"S", // Wool
		"S",
		"S",
		"S",
		"B", // Bruck
		"B",
		"B",
		"O", // Ore
		"O",
		"O",
		"D", // Desert
	}
	rolls := [...]int{
		2,
		3,
		3,
		4,
		4,
		5,
		5,
		6,
		6,
		8,
		8,
		9,
		9,
		10,
		10,
		11,
		11,
		12,
	}
	rand.Seed(time.Now().UnixNano())
	for i := len(tiles) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}
	for i := len(rolls) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		rolls[i], rolls[j] = rolls[j], rolls[i]
	}

	des := 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i] != "D" {
			t := Tile{
				roll:   rolls[i-des],
				bandit: false,
				res:    tiles[i],
			}
			m.tiles = append(m.tiles, &t)
		} else {
			t := Tile{
				roll:   7,
				bandit: true,
				res:    tiles[i],
			}
			m.tiles = append(m.tiles, &t)
			des += 1
		}
	}

	var nodes [][]*Node
	var edges []*Edge
	edgeCounter := 0
	for i := 0; i < len(rowlens); i++ {
		temp := []*Node{}
		if i == 0 { // First row
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}
				n.owner = -1

				var t []*Tile
				if j < rowlens[i]-1 {
					t = append(t, m.tiles[j/2])
				}
				if j > 0 && j%2 == 0 {
					t = append(t, m.tiles[j/2-1])
				}
				sum := 0
				for s := 0; s <= i; s++ {
					sum += rowoff[i]
				}
				n.tiles = t
				n.index = i + sum

				if j == 0 {
					edge1 := Edge{} // Unused
					edge2 := Edge{}
					edge3 := Edge{}

					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 2

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, &edge1, &edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{} // Unused

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge2)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				} else if j%2 == 0 {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 2

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				} else {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{} // Unused
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		} else if i == 5 { // Last row
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}
				n.owner = -1

				var t []*Tile
				if j < rowlens[i]-1 {
					t = append(t, m.tiles[j/2+rowoff[i]])
				}
				if j > 0 && j%2 == 0 {
					t = append(t, m.tiles[j/2+rowoff[i]-1])
				}
				sum := 0
				for s := 0; s <= i; s++ {
					sum += rowoff[i]
				}
				n.tiles = t
				n.index = i + sum

				if j == 0 {
					edge1 := Edge{} // Unused
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, &edge1, edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{} // Unused

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)

					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else if j%2 == 0 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{} // Unused
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		} else if rowlens[i] > rowlens[i-1] { // Row size increasing
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}
				n.owner = -1

				var t []*Tile
				if j == 0 {
					t = append(t, m.tiles[rowoff[i]])
				} else if j == rowlens[i]-1 {
					t = append(t, m.tiles[rowoff[i+1]-1])
				} else if j%2 == 0 {
					t = append(t, m.tiles[j/2+rowoff[i]-1])
					t = append(t, m.tiles[j/2+rowoff[i]])
					t = append(t, m.tiles[j/2+rowoff[i-1]-1])
				} else {
					if j == 1 {
						t = append(t, m.tiles[rowoff[i]])
						t = append(t, m.tiles[rowoff[i-1]])
					} else if j == rowlens[i]-2 {
						t = append(t, m.tiles[rowoff[i+1]-1])
						t = append(t, m.tiles[rowoff[i]-1])
					} else {
						t = append(t, m.tiles[j/2+rowoff[i-1]-1])
						t = append(t, m.tiles[j/2+rowoff[i-1]])
						t = append(t, m.tiles[j/2+rowoff[i]])
					}
				}
				sum := 0
				for s := 0; s <= i; s++ {
					sum += rowoff[i]
				}
				n.tiles = t
				n.index = i + sum

				if j == 0 {
					edge1 := Edge{} // Unused
					edge2 := Edge{}
					edge3 := Edge{}

					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 2

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, &edge1, &edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{} // Unused

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge2)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				} else if j%2 == 0 {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 2
					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				} else {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j-1].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		} else if rowlens[i] < rowlens[i-1] { // Row size decreasing
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}
				n.owner = -1

				var t []*Tile
				if j == 0 {
					t = append(t, m.tiles[rowoff[i]])
				} else if j == rowlens[i]-1 {
					t = append(t, m.tiles[rowoff[i+1]-1])
				} else if j%2 == 0 {
					t = append(t, m.tiles[j/2+rowoff[i]])
					t = append(t, m.tiles[j/2+rowoff[i]-1])
					t = append(t, m.tiles[j/2+rowoff[i+1]-1])
				} else {
					if j == 1 {
						t = append(t, m.tiles[rowoff[i]+1])
						t = append(t, m.tiles[rowoff[i]])
					} else if j == rowlens[i]-2 {
						t = append(t, m.tiles[j/2-1+rowoff[i+1]])
						t = append(t, m.tiles[j/2+rowoff[i]])
					} else {
						t = append(t, m.tiles[j/2+rowoff[i]-1])
						t = append(t, m.tiles[j/2+rowoff[i+1]-1])
						t = append(t, m.tiles[j/2+rowoff[i+1]])
					}
				}
				sum := 0
				for s := 0; s <= i; s++ {
					sum += rowoff[i]
				}
				n.tiles = t
				n.index = i + sum

				if j == 0 {
					edge1 := Edge{} // Unused
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, &edge1, edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{} // Unused

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)

					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else if j%2 == 0 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 2

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		} else { // Row size staying the same
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}
				n.owner = -1

				var t []*Tile
				if j == 0 {
					t = append(t, m.tiles[rowoff[i]])
				} else if j == rowlens[i]-1 {
					t = append(t, m.tiles[rowoff[i+1]-1])
				} else if j%2 == 0 {
					t = append(t, m.tiles[j/2+rowoff[i]-1])
					t = append(t, m.tiles[j/2+rowoff[i]])
					t = append(t, m.tiles[j/2+rowoff[i+1]-1])
				} else {
					if j == 1 {
						t = append(t, m.tiles[rowoff[i]+1])
						t = append(t, m.tiles[rowoff[i]])
					} else if j == rowlens[i]-2 {
						t = append(t, m.tiles[j/2-1+rowoff[i+1]])
						t = append(t, m.tiles[j/2+rowoff[i]])
					} else {
						t = append(t, m.tiles[j/2+rowoff[i]-1])
						t = append(t, m.tiles[j/2+rowoff[i+1]-1])
						t = append(t, m.tiles[j/2+rowoff[i+1]])
					}
				}
				sum := 0
				for s := 0; s <= i; s++ {
					sum += rowoff[i]
				}
				n.tiles = t
				n.index = i + sum

				if j == 0 {
					edge1 := Edge{} // Unused
					edge2 := nodes[i-1][j].edges[1]
					edge3 := Edge{}

					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, &edge1, edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j].edges[1]
					edge3 := Edge{} // Unused

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)

					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else if j%2 == 0 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
					edge2.road = -1
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.road = -1
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 2

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		}
		nodes = append(nodes, temp)
	}
	m.edges = edges
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes[i]); j++ {
			m.nodes = append(m.nodes, nodes[i][j])
		}
	}

	return m
}

func newGame(human bool) Game {
	game := Game{}
	game.board = NewBoard()

	var cards = [...]string{
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"kn",
		"vp",
		"vp",
		"vp",
		"vp",
		"vp",
		"yo",
		"yo",
		"mo",
		"mo",
		"rb",
		"rb",
	}
	rand.Seed(time.Now().UnixNano())
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
	game.cards = cards

	game.players = [...]*Player{
		&Player{},
		&Player{},
		&Player{},
		&Player{},
	}
	if human {
		game.players[rand.Intn(4)].human = true
	}

	return game
}

// We also don't need to talk about this one
// In fact, let's just ignore this whole file
// This absolutely could have been done programmatically, but I was too far in to turn back and offsets are annoying
func PrintGame(game Game) {
	var playerColor = [5]*color.Color{
		color.New(color.FgWhite),
		color.New(color.FgBlue),
		color.New(color.FgRed),
		color.New(color.FgMagenta),
		color.New(color.FgGreen),
	}
	fmt.Println("Player colors are:")
	for i := -1; i < len(playerColor)-1; i++ {
		playerColor[i+1].Printf("Player %d \n", i)
	}

	fmt.Printf("            ")
	playerColor[game.board.nodes[1].owner+1].Printf("*")
	fmt.Printf("       ")
	playerColor[game.board.nodes[3].owner+1].Printf("*")
	fmt.Printf("       ")
	playerColor[game.board.nodes[5].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("          ")
	playerColor[game.board.edges[1].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[2].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[4].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[5].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[7].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[8].road+1].Printf("\\")
	fmt.Printf("\n")

	fmt.Printf("        ")
	playerColor[game.board.nodes[0].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[0].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[2].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[1].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[4].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[2].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[6].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("\n")

	fmt.Printf("        ")
	playerColor[game.board.edges[0].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[3].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[6].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[9].road+1].Printf("|")
	fmt.Printf("\n")

	fmt.Printf("        ")
	playerColor[game.board.nodes[8].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[0].roll)
	playerColor[game.board.nodes[10].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[1].roll)
	playerColor[game.board.nodes[12].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[2].roll)
	playerColor[game.board.nodes[14].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("      ")
	playerColor[game.board.edges[11].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[12].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[14].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[15].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[17].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[18].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[20].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[21].road+1].Printf("\\")
	fmt.Printf("\n")

	fmt.Printf("    ")
	playerColor[game.board.nodes[7].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[3].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[9].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[4].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[11].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[5].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[13].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[6].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[15].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("    ")
	playerColor[game.board.edges[10].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[13].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[16].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[19].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[22].road+1].Printf("|")
	fmt.Printf("\n")

	fmt.Printf("    ")
	playerColor[game.board.nodes[17].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[3].roll)
	playerColor[game.board.nodes[19].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[4].roll)
	playerColor[game.board.nodes[21].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[5].roll)
	playerColor[game.board.nodes[23].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[6].roll)
	playerColor[game.board.nodes[25].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("  ")
	playerColor[game.board.edges[24].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[25].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[27].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[28].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[30].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[31].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[33].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[34].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[36].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[37].road+1].Printf("\\")
	fmt.Printf("\n")

	playerColor[game.board.nodes[16].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[7].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[18].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[8].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[20].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[9].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[22].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[10].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[24].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[11].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[26].owner+1].Printf("*")
	fmt.Printf("\n")

	playerColor[game.board.edges[23].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[26].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[29].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[32].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[35].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[38].road+1].Printf("|")
	fmt.Printf("\n")

	playerColor[game.board.nodes[27].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[7].roll)
	playerColor[game.board.nodes[29].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[8].roll)
	playerColor[game.board.nodes[31].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[9].roll)
	playerColor[game.board.nodes[33].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[10].roll)
	playerColor[game.board.nodes[35].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[11].roll)
	playerColor[game.board.nodes[37].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("  ")
	playerColor[game.board.edges[39].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[41].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[42].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[44].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[45].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[47].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[48].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[50].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[51].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[53].road+1].Printf("/")
	fmt.Printf("\n")

	fmt.Printf("    ")
	playerColor[game.board.nodes[28].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[12].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[30].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[13].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[32].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[14].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[34].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[15].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[36].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("    ")
	playerColor[game.board.edges[40].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[43].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[46].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[49].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[52].road+1].Printf("|")
	fmt.Printf("\n")

	fmt.Printf("    ")
	playerColor[game.board.nodes[38].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[12].roll)
	playerColor[game.board.nodes[40].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[13].roll)
	playerColor[game.board.nodes[42].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[14].roll)
	playerColor[game.board.nodes[44].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[15].roll)
	playerColor[game.board.nodes[46].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("      ")
	playerColor[game.board.edges[54].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[56].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[57].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[59].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[60].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[62].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[63].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[65].road+1].Printf("/")
	fmt.Printf("\n")

	fmt.Printf("        ")
	playerColor[game.board.nodes[39].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[16].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[41].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[17].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[43].owner+1].Printf("*")
	fmt.Printf("   ")
	fmt.Printf("%s", game.board.tiles[18].res)
	fmt.Printf("   ")
	playerColor[game.board.nodes[45].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("        ")
	playerColor[game.board.edges[55].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[58].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[61].road+1].Printf("|")
	fmt.Printf("       ")
	playerColor[game.board.edges[64].road+1].Printf("|")
	fmt.Printf("\n")

	fmt.Printf("        ")
	playerColor[game.board.nodes[47].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[16].roll)
	playerColor[game.board.nodes[49].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[17].roll)
	playerColor[game.board.nodes[51].owner+1].Printf("*")
	fmt.Printf("  %2d   ", game.board.tiles[18].roll)
	playerColor[game.board.nodes[53].owner+1].Printf("*")
	fmt.Printf("\n")

	fmt.Printf("          ")
	playerColor[game.board.edges[66].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[67].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[68].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[69].road+1].Printf("/")
	fmt.Printf("   ")
	playerColor[game.board.edges[70].road+1].Printf("\\")
	fmt.Printf("   ")
	playerColor[game.board.edges[71].road+1].Printf("/")
	fmt.Printf("\n")

	fmt.Printf("            ")
	playerColor[game.board.nodes[48].owner+1].Printf("*")
	fmt.Printf("       ")
	playerColor[game.board.nodes[50].owner+1].Printf("*")
	fmt.Printf("       ")
	playerColor[game.board.nodes[52].owner+1].Printf("*")
	fmt.Printf("\n")
}
