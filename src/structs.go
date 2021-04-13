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

					fmt.Println(edge2.index)

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
	var playerColor = [4]*color.Color{
		color.New(color.FgBlue),
		color.New(color.FgRed),
		color.New(color.FgMagenta),
		color.New(color.FgGreen),
	}
	none := color.New(color.FgWhite)
	fmt.Println("Player colors are:")
	none.Println("No Player")
	for i := 0; i < len(playerColor); i++ {
		playerColor[i].Printf("Player %d \n", i+1)
	}

}
