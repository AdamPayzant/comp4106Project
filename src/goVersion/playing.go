package main

import (
	"sync"
	"sync/atomic"
)

const (
	maxdepth int = 5
)

func main() {
	// Game setup
	game := newGame()
	var i int
	for i = 0; i < len(game.players); i++ {
		placeStart(&game.board, game.players[i])
	}
	for ; i >= 0; i-- {
		placeStart(&game.board, game.players[i])
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
		for i := 0; i < len(player.roads); i++ {
			for j := 0; j < 3; j++ {
				if player.roads[i].nodes[0].edges[j].inUse && player.roads[i].nodes[0].edges[j].road == 0 {
					m := Move{}
					moves = append(moves, &m)

					m.cost["cl"] += 1
					m.cost["wo"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[0].edges[j].index,
						nodes: player.roads[i].nodes[0].edges[j].nodes,
						road:  player.number,
					})
					wg.Add(1)
					go continueMove(&wg, player, game, &m, 0)
				}
				if player.roads[i].nodes[1].edges[j].inUse && player.roads[i].nodes[1].edges[j].road == 0 {
					m := Move{}
					moves = append(moves, &m)

					m.cost["cl"] += 1
					m.cost["wo"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[1].edges[j].index,
						nodes: player.roads[i].nodes[1].edges[j].nodes,
						road:  player.number,
					})
					wg.Add(1)
					go continueMove(&wg, player, game, &m, 0)
				}
			}
		}
	}
	// Settlement case
	if player.res["cl"] >= 1 && player.res["wh"] >= 1 && player.res["sh"] >= 1 && player.res["wo"] >= 1 {
		for i := 0; i < len(player.roads); i++ {
			if isBuildable(player.roads[i].nodes[0]) {
				m := Move{}
				moves = append(moves, &m)

				m.cost["cl"] += 1
				m.cost["wh"] += 1
				m.cost["sh"] += 1
				m.cost["wo"] += 1
				m.newNodes = append(m.newNodes, Node{
					index:    player.roads[i].nodes[0].index,
					building: 1,
					owner:    player.number,
					edges:    player.roads[i].nodes[0].edges,
					tiles:    player.roads[i].nodes[0].tiles,
				})
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
			if isBuildable(player.roads[i].nodes[1]) {
				m := Move{}
				moves = append(moves, &m)

				m.cost["cl"] += 1
				m.cost["wh"] += 1
				m.cost["sh"] += 1
				m.cost["wo"] += 1
				m.newNodes = append(m.newNodes, Node{
					index:    player.roads[i].nodes[1].index,
					building: 1,
					owner:    player.number,
					edges:    player.roads[i].nodes[1].edges,
					tiles:    player.roads[i].nodes[1].tiles,
				})
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
		}
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
			go continueMove(&wg, player, game, &m, 0)
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

func continueMove(wg *sync.WaitGroup, player *Player, game Game, move *Move, depth int) {
	defer wg.Done()
	var moves []*Move
	var wg2 sync.WaitGroup
	// Road case
	if player.res["cl"]-move.cost["cl"] >= 1 && player.res["wo"]-move.cost["wo"] >= 1 {
		for i := 0; i < len(player.roads); i++ {
			for j := 0; j < 3; j++ {
				found1 := false
				found2 := false
				for n := 0; n < len(move.newEdges); n++ {
					if player.roads[i].nodes[0].edges[j].index == move.newEdges[n].index {
						found1 = true
					}
					if player.roads[i].nodes[1].edges[j].index == move.newEdges[n].index {
						found2 = true
					}
				}
				if player.roads[i].nodes[0].edges[j].inUse && player.roads[i].nodes[0].edges[j].road == 0 && !found1 {
					m := Move{}
					moves = append(moves, &m)

					m.cost["cl"] += 1
					m.cost["wo"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[0].edges[j].index,
						nodes: player.roads[i].nodes[0].edges[j].nodes,
						road:  player.number,
					})
					wg.Add(1)
					go continueMove(&wg2, player, game, &m, depth)
				}
				if player.roads[i].nodes[1].edges[j].inUse && player.roads[i].nodes[1].edges[j].road == 0 && !found2 {
					m := Move{}
					moves = append(moves, &m)

					m.cost["cl"] += 1
					m.cost["wo"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[1].edges[j].index,
						nodes: player.roads[i].nodes[1].edges[j].nodes,
						road:  player.number,
					})
					wg.Add(1)
					go continueMove(&wg2, player, game, &m, depth)
				}
			}
		}
	}
	// Settlement case
	if player.res["cl"]-move.cost["cl"] >= 1 && player.res["wh"]-move.cost["wh"] >= 1 &&
		player.res["sh"]-move.cost["sh"] >= 1 && player.res["wo"]-move.cost["wo"] >= 1 {
		for i := 0; i < len(player.roads); i++ {
			found1 := false
			found2 := false
			for j := 0; j < len(move.newNodes); j++ {
				if player.roads[i].nodes[0].index == move.newNodes[j].index {
					found1 = true
				}
				if player.roads[i].nodes[0].index == move.newNodes[j].index {
					found2 = true
				}
			}
			if isBuildable(player.roads[i].nodes[0]) && !found1 {
				m := Move{}
				moves = append(moves, &m)

				m.cost["cl"] += 1
				m.cost["wh"] += 1
				m.cost["sh"] += 1
				m.cost["wo"] += 1
				m.newNodes = append(m.newNodes, Node{
					index:    player.roads[i].nodes[0].index,
					building: 1,
					owner:    player.number,
					edges:    player.roads[i].nodes[0].edges,
					tiles:    player.roads[i].nodes[0].tiles,
				})
				wg.Add(1)
				go continueMove(&wg2, player, game, &m, depth)
			}
			if isBuildable(player.roads[i].nodes[1]) && !found2 {
				m := Move{}
				moves = append(moves, &m)

				m.cost["cl"] += 1
				m.cost["wh"] += 1
				m.cost["sh"] += 1
				m.cost["wo"] += 1
				m.newNodes = append(m.newNodes, Node{
					index:    player.roads[i].nodes[1].index,
					building: 1,
					owner:    player.number,
					edges:    player.roads[i].nodes[1].edges,
					tiles:    player.roads[i].nodes[1].tiles,
				})
				wg.Add(1)
				go continueMove(&wg2, player, game, &m, depth)
			}
		}
	}
	// City case
	if player.res["wh"]-move.cost["wh"] >= 2 && player.res["ir"]-move.cost["ir"] >= 3 {
		for i := 0; i < len(player.villages); i++ {
			found := false
			for j := 0; j < len(move.newNodes); j++ {
				if move.newNodes[j].index == player.villages[i].index && move.newNodes[j].building == 2 {
					found = true
					break
				}
			}
			if !found {
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
				go continueMove(&wg2, player, game, &m, depth)
			}
		}
	}
	// Pass
	m := Move{}
	moves = append(moves, &m)
	nextPlayer := player.number + 1
	if nextPlayer > 3 {
		nextPlayer = 0
	}
	// TODO: DO COMPLETE COPY OF GAME INTO newGAME
	newGame := game
	go predict(&wg2, *game.players[nextPlayer], player.number, newGame, depth, &m.heur)
}

func humanPlay() Move {
	return Move{}
}

func predict(wg *sync.WaitGroup, player Player, hostID int, game Game, depth int, res *float32) {
	defer wg.Done()
	if depth == maxdepth {
		*res = 0
		return
	}
}

func playMove(move Move, p Player, g Game) Game {
	game := g
	player := p

	if move.card != "" {
		m := map[string]func(*Game, *Player){
			"kn": knight,
			"mo": monopoly,
			"yo": yop,
			"rb": roadBuilding,
		}
		m[move.card](&game, &player)
	}

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
	// TODO: Check if the new path creates a new longest road
	game.players[player.number-1] = &player

	return game
}

func placeStart(board *Board, player *Player) {
	var wg sync.WaitGroup
	var heur []float32
	for i := 0; i < len(board.nodes); i++ {
		var f float32
		heur = append(heur, f)
		wg.Add(1)
		go settleHeur(&wg, board, board.nodes[i], player, &f)
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
	player.villages = append(player.villages, board.nodes[bestIn])
	if len(player.villages) > 1 {
		if player.villages[0].index > player.villages[1].index {

		} else {

		}
	} else {

	}
}

func isBuildable(node *Node) bool {
	if node.building != 0 {
		return false
	}
	for i := 0; i < len(node.edges); i++ {
		if node.edges[i].nodes[0].index != node.index {
			if node.edges[i].nodes[0].building != 0 {
				return false
			}
		} else {
			if node.edges[i].nodes[1].building != 0 {
				return false
			}
		}
	}
	return true
}

func settleHeur(wg *sync.WaitGroup, board *Board, node *Node, player *Player, heur *float32) {
	defer wg.Done()
	if !isBuildable(node) {
		*heur = 0
	}
	var count int32
	var explored map[int]bool
	var wg2 sync.WaitGroup
	var mut sync.Mutex
	wg2.Add(1)
	searchVil(&wg2, &mut, node, 4, &count, &explored)
	wg2.Wait()
	var res map[string]float32
	for i := 0; i < len(node.tiles); i++ {
		res[node.tiles[i].res] += diceOdds[node.tiles[i].roll-1]
	}
	for i := 0; i < len(player.villages); i++ {
		for j := 0; j < len(player.villages[i].tiles); j++ {
			res[player.villages[i].tiles[j].res] += diceOdds[player.villages[i].tiles[j].roll-1]
		}
	}

	var totalOdds float32 = 0
	for _, val := range res {
		totalOdds += val
	}

	*heur = (totalOdds * 10) + (float32(len(res)) * 2) - float32(count)
}

func searchVil(wg *sync.WaitGroup, mut *sync.Mutex, node *Node, depth int, count *int32, explored *map[int]bool) {
	defer wg.Done()
	// This lock kills performance
	// I really don't want this but there's not much I can do while minimizing space complexity
	mut.Lock()
	if _, ok := (*explored)[node.index]; !ok {
		mut.Unlock()
		return
	}
	(*explored)[node.index] = true
	mut.Unlock()
	var wg2 sync.WaitGroup
	if node.building > 0 {
		atomic.AddInt32(count, 1)
	}
	if depth != 0 {
		for i := 0; i < 3; i++ {
			if node.edges[i].inUse {
				wg2.Add(1)
				if node.edges[i].nodes[0].index != node.index {
					go searchVil(&wg2, mut, node.edges[i].nodes[0], depth-1, count, explored)
				} else {
					go searchVil(&wg2, mut, node.edges[i].nodes[1], depth-1, count, explored)
				}
			}
		}
	}
	wg2.Wait()
}
