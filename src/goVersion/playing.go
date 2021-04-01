package main

import (
	"sync"
)

const (
	maxdepth int = 5
)

func main() {
	// Game setup
	game := newGame()
	var i int
	for i = 0; i < len(game.players); i++ {
		placeStart(&game.board, *game.players[i])
	}
	for ; i >= 0; i-- {
		placeStart(&game.board, *game.players[i])
	}
	// Game loop
	curInd := 3
	for game.players[curInd].victoryPoints < 10 {
		if curInd == 3 {
			curInd = 0
		} else {
			curInd++
		}

		move := play(game.players[curInd], game)
		game = playMove(move, *game.players[curInd], game)
	}
}

func play(player *Player, game Game) Move {
	if player.human {
		return humanPlay()
	}
	var moves []*Move
	var wg sync.WaitGroup
	// Road case
	if player.res["cl"] >= 1 && player.res["wo"] >= 1 {

	}
	// Settlement case
	if player.res["cl"] >= 1 && player.res["wh"] >= 1 && player.res["sh"] >= 1 && player.res["wo"] >= 1 {

	}
	// City case
	if player.res["wh"] >= 2 && player.res["ir"] >= 3 {
		for i := 0; i < len(player.villages); i++ {
			m := Move{}
			moves = append(moves, &m)

			m.cost["wh"] = 2
			m.cost["ir"] = 3
			m.newNodes = append(m.newNodes, Node{
				index:    player.villages[i].index,
				building: 2,
				owner:    player.villages[i].owner,
				edges:    player.villages[i].edges,
				tiles:    player.villages[i].tiles,
			})
			wg.Add(1)
			go continueMove(&wg, player, game, &m)
		}
	}
	// Card case
	if player.res["sh"] >= 1 && player.res["wh"] >= 1 && player.res["ir"] >= 1 {

	}
	// Pass
	m := Move{}
	moves = append(moves, &m)
	nextPlayer := player.number + 1
	if nextPlayer > 3 {
		nextPlayer = 0
	}
	go predict(&wg, *game.players[nextPlayer], player.number, game, 0, &m.heur)

	wg.Wait()
	bestVal := moves[0].heur
	bestInd := 0
	for i := 1; i < len(moves); i++ {
		if moves[i].heur > bestVal {
			bestVal = moves[i].heur
			bestInd = i
		}
	}
	return *moves[bestInd]
}

func continueMove(wg *sync.WaitGroup, player *Player, game Game, move *Move) {

}

func humanPlay() Move {
	return Move{}
}

func predict(wg *sync.WaitGroup, player Player, hostID int, game Game, depth int, res *float32) {
	wg.Done()
}

func playMove(move Move, p Player, g Game) Game {
	game := g
	player := p

	for k, v := range move.cost {
		player.res[k] -= v
	}
	for i := 0; i < len(move.newNodes); i++ {
		for j := 0; j < 3; j++ {
			if game.board.nodes[i].edges[j].inUse {
				if game.board.nodes[move.newNodes[i].index].edges[j].nodes[0].index == move.newNodes[i].index {
					game.board.nodes[move.newNodes[i].index].edges[j].nodes[0] = &move.newNodes[i]
				} else {
					game.board.nodes[move.newNodes[i].index].edges[j].nodes[1] = &move.newNodes[i]
				}
			}
		}
		game.board.nodes[i] = &move.newNodes[i]
	}
	for i := 0; i < len(move.newEdges); i++ {
		move.newEdges[i].nodes = game.board.nodes
		game.board.edges[i] = &move.newEdges[i]
	}

	game.players[player.number-1] = &player

	return game
}

// TODO
func placeStart(board *Board, player Player) {
	var wg sync.WaitGroup
	var heur []float32
	for i := 0; i < len(board.nodes); i++ {
		var f float32
		heur = append(heur, f)
		go startHeur(&wg, board, board.nodes[i], player.number, &f)
	}

	var max float32 = 0
	var bestIn int
	for i := 0; i < len(heur); i++ {
		if heur[i] > max {
			max = heur[i]
			bestIn = i
		}
	}
	board.nodes[bestIn].building = 1
	board.nodes[bestIn].owner = player.number
}

func isBuildable(node *Node) bool {
	return true
}

func startHeur(wg *sync.WaitGroup, board *Board, node *Node, playerNum int, heur *float32) {
	if !isBuildable(node) {
		*heur = 0
	}

}
