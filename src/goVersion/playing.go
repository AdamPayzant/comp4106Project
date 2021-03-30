package main

import "sync"

const (
	maxdepth int = 5
)

func main() {
	// TODO: Do setup
	// TODO: Peform game loop
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

}
