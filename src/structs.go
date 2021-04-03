package main

import (
	"math/rand"
	"time"
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
	players []*Player
	board   Board

	longestRoad int
	lrPlayer    int
	largestArmy int
	laPlayer    int
}

func NewBoard() Board {
	m := Board{}
	tiles := [...]string{
		"wo",
		"wo",
		"wo",
		"wo",
		"wh",
		"wh",
		"wh",
		"wh",
		"sh",
		"sh",
		"sh",
		"sh",
		"cl",
		"cl",
		"cl",
		"ir",
		"ir",
		"ir",
		"de",
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
	for i := 0; i < len(rolls); i++ {
		if tiles[i] != "de" {
			t := Tile{
				roll:   rolls[i],
				bandit: false,
				res:    tiles[i+des],
			}
			m.tiles = append(m.tiles, &t)
		}
	}

	var nodes [][]*Node
	var edges []*Edge
	edgeCounter := 0
	for i := 0; i < len(rowlens); i++ {
		temp := make([]*Node, rowlens[i])
		if i == 0 { // First row
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}

				var t []*Tile
				if j < rowlens[i]-1 {
					t = append(t, m.tiles[j/2+rowoff[i]])
				}
				if j > 0 && i%2 == 0 {
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
					edge2 := Edge{}
					edge3 := Edge{}

					edge2.inUse = true
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
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
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
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

				var t []*Tile
				if j < rowlens[i]-1 {
					t = append(t, m.tiles[j/2+rowoff[i]])
				}
				if j > 0 && i%2 == 0 {
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
					edge3.index = edgeCounter + 1
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
					edge2 := Edge{} // Unused
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge3.inUse = true
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				} else {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		} else if rowlens[i] > rowlens[i-1] { // Row size increasing
			for j := 0; j < rowlens[i]; j++ {
				n := Node{}

				var t []*Tile
				if j == 0 {
					t = append(t, m.tiles[rowoff[i]])
				} else if j == rowlens[i]-1 {
					t = append(t, m.tiles[j/2+rowoff[i]-1])
				} else if j%2 == 0 {
					t = append(t, m.tiles[j/2+rowoff[i]-1])
					t = append(t, m.tiles[j/2+rowoff[i]])
					t = append(t, m.tiles[j/2+rowoff[i-1]-1])
				} else {
					if j == 1 {
						t = append(t, m.tiles[rowoff[i]])
						t = append(t, m.tiles[rowoff[i-1]])
					} else if j == rowlens[i]-2 {
						t = append(t, m.tiles[j/2+rowoff[i]])
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
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 3

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, &edge1, &edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := Edge{}
					edge3 := Edge{} // Unused

					edge1.nodes = append(edge1.nodes, &n)
					edge2.inUse = true
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
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
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
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, &edge1, edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, edge1, edge2, &edge3)
				} else if j%2 == 0 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j+1].edges[1]
					edge3 := Edge{}

					edge1.nodes = append(edge1.nodes, &n)
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
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
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		} else {
			for j := 0; j < rowlens[i]; j++ { // Row size staying the same
				n := Node{}

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
					edge3.index = edgeCounter
					edge3.nodes = append(edge3.nodes, &n)
					edgeCounter += 1

					edges = append(edges, &edge3)
					n.edges = append(n.edges, &edge1, edge2, &edge3)
				} else if j == rowlens[i]-1 {
					edge1 := temp[j-1].edges[2]
					edge2 := nodes[i-1][j].edges[1]
					edge3 := Edge{}

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
					edge2.index = edgeCounter
					edge2.nodes = append(edge2.nodes, &n)
					edge3.inUse = true
					edge3.index = edgeCounter + 1
					edge3.nodes = append(edge3.nodes, &n)

					edges = append(edges, &edge2, &edge3)
					n.edges = append(n.edges, edge1, &edge2, &edge3)
				}
				temp = append(temp, &n)
			}
		}
		nodes = append(nodes, temp)
	}

	return m
}

func newGame() Game {
	return Game{}
}

func PrintBoard(board Board) {

}
