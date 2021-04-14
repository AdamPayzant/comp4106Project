package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxdepth     int = 0
	maxMoveDepth int = 2
	// Weights
	resVariety float32 = 10
	resOdds    float32 = 10
	lrDistance float32 = 2
	laDistance float32 = 1
	vp         float32 = 1
	resInHand  float32 = 1
)

func main() {
	// Game setup
	game := newGame(false)

	var i int
	for i = 0; i < len(game.players); i++ {
		placeStart(&game.board, game.players[i])
	}
	i--
	for ; i >= 0; i-- {
		placeStart(&game.board, game.players[i])
	}
	PrintGame(game)

	for i := 0; i < len(game.board.nodes); i++ {
		fmt.Println(isBuildable(game.board.nodes[i]))
	}

	// Game loop
	curInd := 3
	//for i := 0; i < 6; i++ {
	for game.players[curInd].victoryPoints < 10 {
		if curInd == 3 {
			curInd = 0
		} else {
			curInd++
		}
		roll := rand.Intn(11) + 2
		fmt.Println(roll)
		if roll != 7 {
			distRes(&game, roll)
		}

		fmt.Printf("Player %d's turn \n", curInd)
		move := play(game.players[curInd], game)
		game = playMove(move, curInd, game)
		PrintGame(game)
	}
	fmt.Printf("PLAYER %d WON!\n", curInd)
}

func distRes(game *Game, roll int) {
	for i := 0; i < len(game.players); i++ {
		for j := 0; j < len(game.players[i].villages); j++ {
			for k := 0; k < len(game.players[i].villages[j].tiles); k++ {
				if !game.players[i].villages[j].tiles[k].bandit && game.players[i].villages[j].tiles[k].roll == roll {
					game.players[i].res[game.players[i].villages[j].tiles[k].res] += game.players[i].villages[j].building
				}
			}
		}
	}
}

func play(player *Player, game Game) Move {
	fmt.Println(player.res)
	if player.human {
		return humanPlay()
	}
	var moves []*Move
	var wg sync.WaitGroup
	// Road case
	if player.res["B"] >= 1 && player.res["L"] >= 1 {
		fmt.Println("ROAD CASE")
		for i := 0; i < len(player.roads); i++ {
			for j := 0; j < 3; j++ {
				if player.roads[i].nodes[0].edges[j].inUse &&
					player.roads[i].nodes[0].edges[j].road == -1 {
					m := newMove()
					moves = append(moves, &m)

					m.cost["B"] += 1
					m.cost["L"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[0].edges[j].index,
						nodes: player.roads[i].nodes[0].edges[j].nodes,
						road:  player.number,
					})
					wg.Add(1)
					go continueMove(&wg, player, game, &m, 0)
				}
				if player.roads[i].nodes[1].edges[j].inUse && player.roads[i].nodes[1].edges[j].road == -1 {
					m := newMove()
					moves = append(moves, &m)

					m.cost["B"] += 1
					m.cost["L"] += 1
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
	if player.res["B"] >= 1 && player.res["W"] >= 1 && player.res["S"] >= 1 && player.res["L"] >= 1 {
		fmt.Println("SETTLEMENT CASE")
		for i := 0; i < len(player.roads); i++ {
			if isBuildable(player.roads[i].nodes[0]) {
				fmt.Printf("Build site found")
				m := newMove()
				moves = append(moves, &m)

				m.cost["B"] += 1
				m.cost["W"] += 1
				m.cost["S"] += 1
				m.cost["L"] += 1
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
				fmt.Printf("Build site found")
				m := newMove()
				moves = append(moves, &m)

				m.cost["B"] += 1
				m.cost["W"] += 1
				m.cost["S"] += 1
				m.cost["L"] += 1
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
	if player.res["W"] >= 2 && player.res["O"] >= 3 {
		fmt.Println("CITY CASE")
		for i := 0; i < len(player.villages); i++ {
			if player.villages[i].building == 1 {
				m := newMove()
				moves = append(moves, &m)

				m.cost["W"] += 2
				m.cost["O"] += 3
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
	}
	// Buy Card case
	if player.res["S"] >= 1 && player.res["W"] >= 1 && player.res["O"] >= 1 {
		// TODO
	}
	// Play card case
	for i := 0; i < len(player.cards); i++ {
		m := newMove()
		m.card = player.cards[i]
		wg.Add(1)
		g := playMove(m, player.number, game)
		wg.Add(1)
		go continueMove(&wg, g.players[player.number], g, &m, 0)
	}
	// World trade case
	for k, val := range player.res {
		if val > 4 {
			if player.res["O"] < 3 {
				m := newMove()
				m.cost[k] += 4

				m.gain["O"] += 1
				moves = append(moves, &m)
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
			if player.res["B"] < 3 {
				m := newMove()
				m.cost[k] += 4

				m.gain["B"] += 1
				moves = append(moves, &m)
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
			if player.res["S"] < 3 {
				m := newMove()
				m.cost[k] += 4

				m.gain["S"] += 1
				moves = append(moves, &m)
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
			if player.res["W"] < 3 {
				m := newMove()
				m.cost[k] += 4

				m.gain["W"] += 1
				moves = append(moves, &m)
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
			if player.res["L"] < 3 {
				m := newMove()
				m.cost[k] += 4

				m.gain["L"] += 1
				moves = append(moves, &m)
				wg.Add(1)
				go continueMove(&wg, player, game, &m, 0)
			}
		}
	}
	// Pass
	m := newMove()
	m.heur = heuristic(game, player.number) * (1 / 10)
	moves = append(moves, &m)
	nextPlayer := player.number + 1
	if nextPlayer > 3 {
		nextPlayer = 0
	}
	wg.Add(1)
	go predict(&wg, *game.players[nextPlayer], player.number, game, 0, &m.heur)

	fmt.Println("Waiting...")
	wg.Wait()
	bestInd := 0
	for i := 1; i < len(moves); i++ {
		if moves[i].heur > moves[bestInd].heur {
			bestInd = i
		}
	}
	fmt.Println(*moves[bestInd])
	return *moves[bestInd]
}

func continueMove(wg *sync.WaitGroup, player *Player, game Game, move *Move, depth int) {
	defer wg.Done()

	var moves []*Move
	var wg2 sync.WaitGroup
	if move.moveDepth < maxMoveDepth {
		// Road case
		if player.res["B"]-move.cost["B"]+move.gain["B"] >= 1 &&
			player.res["L"]-move.cost["L"]+move.gain["L"] >= 1 {
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
					if player.roads[i].nodes[0].edges[j].inUse &&
						player.roads[i].nodes[0].edges[j].road == -1 && !found1 {
						m := copyMove(*move)
						moves = append(moves, &m)

						m.cost["B"] += 1
						m.cost["L"] += 1
						m.newEdges = append(m.newEdges, Edge{
							inUse: true,
							index: player.roads[i].nodes[0].edges[j].index,
							nodes: player.roads[i].nodes[0].edges[j].nodes,
							road:  player.number,
						})
						wg2.Add(1)
						go continueMove(&wg2, player, game, &m, depth)
					}
					if player.roads[i].nodes[1].edges[j].inUse &&
						player.roads[i].nodes[1].edges[j].road == -1 && !found2 {
						m := copyMove(*move)
						moves = append(moves, &m)

						m.cost["B"] += 1
						m.cost["L"] += 1
						m.newEdges = append(m.newEdges, Edge{
							inUse: true,
							index: player.roads[i].nodes[1].edges[j].index,
							nodes: player.roads[i].nodes[1].edges[j].nodes,
							road:  player.number,
						})
						wg2.Add(1)
						go continueMove(&wg2, player, game, &m, depth)
					}
				}
			}
		}
		// Settlement case
		if player.res["B"]-move.cost["B"]+move.gain["B"] >= 1 &&
			player.res["W"]-move.cost["W"]+move.gain["W"] >= 1 &&
			player.res["S"]-move.cost["S"] >= 1 && player.res["L"]-move.cost["L"] >= 1 {
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
					m := copyMove(*move)
					moves = append(moves, &m)

					m.cost["B"] += 1
					m.cost["W"] += 1
					m.cost["S"] += 1
					m.cost["L"] += 1
					m.newNodes = append(m.newNodes, Node{
						index:    player.roads[i].nodes[0].index,
						building: 1,
						owner:    player.number,
						edges:    player.roads[i].nodes[0].edges,
						tiles:    player.roads[i].nodes[0].tiles,
					})
					wg2.Add(1)
					go continueMove(&wg2, player, game, &m, depth)
				}
				if isBuildable(player.roads[i].nodes[1]) && !found2 {
					m := newMove()
					moves = append(moves, &m)

					m.cost["B"] += 1
					m.cost["W"] += 1
					m.cost["S"] += 1
					m.cost["L"] += 1
					m.newNodes = append(m.newNodes, Node{
						index:    player.roads[i].nodes[1].index,
						building: 1,
						owner:    player.number,
						edges:    player.roads[i].nodes[1].edges,
						tiles:    player.roads[i].nodes[1].tiles,
					})
					wg2.Add(1)
					go continueMove(&wg2, player, game, &m, depth)
				}
			}
		}
		// City case
		if player.res["W"]-move.cost["W"]+move.gain["W"] >= 2 &&
			player.res["O"]-move.cost["O"]+move.gain["O"] >= 3 {
			for i := 0; i < len(player.villages); i++ {
				found := false
				for j := 0; j < len(move.newNodes); j++ {
					if move.newNodes[j].index == player.villages[i].index && move.newNodes[j].building == 2 {
						found = true
						break
					}
				}
				if !found && player.villages[i].building != 2 {
					m := copyMove(*move)
					moves = append(moves, &m)

					m.cost["W"] += 2
					m.cost["O"] += 3
					m.newNodes = append(m.newNodes, Node{
						index:    player.villages[i].index,
						building: 2,
						owner:    player.villages[i].owner,
						edges:    player.villages[i].edges,
						tiles:    player.villages[i].tiles,
					})
					wg2.Add(1)
					go continueMove(&wg2, player, game, &m, depth)
				}
			}
		}
	}
	// Pass
	m := copyMove(*move)
	moves = append(moves, &m)
	nextPlayer := player.number * (depth / 10)
	if nextPlayer > 3 {
		nextPlayer = 0
	}
	newGame := playMove(*move, player.number, game)

	m.heur = heuristic(newGame, player.number)
	wg2.Add(1)
	go predict(&wg2, *game.players[nextPlayer], player.number, newGame, depth, &m.heur)

	wg2.Wait()
	greatestind := 0
	for i := 0; i < len(moves); i++ {
		if moves[i].heur > moves[greatestind].heur {
			greatestind = i
		}
	}

	*move = *moves[greatestind]
}

func humanPlay() Move {
	return newMove()
}

func predict(wg *sync.WaitGroup, player Player, hostID int, game Game, depth int, res *float32) {
	defer wg.Done()
	if depth == maxdepth {
		return
	}

	var moves []*Move
	var wg2 sync.WaitGroup
	// Road case
	if player.res["B"] >= 1 && player.res["L"] >= 1 {
		for i := 0; i < len(player.roads); i++ {
			for j := 0; j < 3; j++ {
				if player.roads[i].nodes[0].edges[j].inUse &&
					player.roads[i].nodes[0].edges[j].road == -1 {
					m := newMove()
					moves = append(moves, &m)

					m.cost["B"] += 1
					m.cost["L"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[0].edges[j].index,
						nodes: player.roads[i].nodes[0].edges[j].nodes,
						road:  player.number,
					})
					wg2.Add(1)
					go continueMove(&wg2, &player, game, &m, depth)
				}
				if player.roads[i].nodes[1].edges[j].inUse &&
					player.roads[i].nodes[1].edges[j].road == -1 {
					m := newMove()
					moves = append(moves, &m)

					m.cost["B"] += 1
					m.cost["L"] += 1
					m.newEdges = append(m.newEdges, Edge{
						inUse: true,
						index: player.roads[i].nodes[1].edges[j].index,
						nodes: player.roads[i].nodes[1].edges[j].nodes,
						road:  player.number,
					})
					wg2.Add(1)
					go continueMove(&wg2, &player, game, &m, depth)
				}
			}
		}
	}
	// Settlement case
	if player.res["B"] >= 1 && player.res["W"] >= 1 && player.res["S"] >= 1 && player.res["L"] >= 1 {
		for i := 0; i < len(player.roads); i++ {
			if isBuildable(player.roads[i].nodes[0]) {
				m := newMove()
				moves = append(moves, &m)

				m.cost["B"] += 1
				m.cost["W"] += 1
				m.cost["S"] += 1
				m.cost["L"] += 1
				m.newNodes = append(m.newNodes, Node{
					index:    player.roads[i].nodes[0].index,
					building: 1,
					owner:    player.number,
					edges:    player.roads[i].nodes[0].edges,
					tiles:    player.roads[i].nodes[0].tiles,
				})
				wg2.Add(1)
				go continueMove(&wg2, &player, game, &m, depth)
			}
			if isBuildable(player.roads[i].nodes[1]) {
				m := newMove()
				moves = append(moves, &m)

				m.cost["B"] += 1
				m.cost["W"] += 1
				m.cost["S"] += 1
				m.cost["L"] += 1
				m.newNodes = append(m.newNodes, Node{
					index:    player.roads[i].nodes[1].index,
					building: 1,
					owner:    player.number,
					edges:    player.roads[i].nodes[1].edges,
					tiles:    player.roads[i].nodes[1].tiles,
				})
				wg2.Add(1)
				go continueMove(&wg2, &player, game, &m, depth)
			}
		}
	}
	// City case
	if player.res["W"] >= 2 && player.res["O"] >= 3 {
		for i := 0; i < len(player.villages); i++ {
			if player.villages[i].building == 1 {
				m := newMove()
				moves = append(moves, &m)

				m.cost["W"] += 2
				m.cost["O"] += 3
				m.newNodes = append(m.newNodes, Node{
					index:    player.villages[i].index,
					building: 2,
					owner:    player.villages[i].owner,
					edges:    player.villages[i].edges,
					tiles:    player.villages[i].tiles,
				})
				wg2.Add(1)
				go continueMove(&wg2, &player, game, &m, depth)
			}
		}
	}
	// Pass
	m := newMove()
	m.heur = heuristic(game, player.number) * float32(depth/10)
	moves = append(moves, &m)
	nextPlayer := player.number + 1
	if nextPlayer > 3 {
		nextPlayer = 0
	}
	wg2.Add(1)
	go predict(&wg2, *game.players[nextPlayer], player.number, game, depth+1, &m.heur)

	wg2.Wait()
	bestVal := moves[0].heur
	bestInd := 0
	for i := 1; i < len(moves); i++ {
		if moves[i].heur > bestVal {
			bestVal = moves[i].heur
			bestInd = i
		}
	}
	h := heuristic(playMove(*moves[bestInd], player.number, game), player.number)
	*res -= h
	return
}

func heuristic(game Game, playerNum int) float32 {
	heur := float32(0)
	player := game.players[playerNum]

	// Both resource odds and variety
	res := make(map[string]float32)
	for i := 0; i < len(player.villages); i++ {
		for j := 0; j < len(player.villages[i].tiles); j++ {
			res[player.villages[i].tiles[j].res] += float32(player.villages[i].building) * diceOdds[player.villages[i].tiles[j].roll-1]
		}
	}
	heur += float32(len(res)) / 5 * resVariety
	var sum float32
	for _, val := range res {
		sum += val
	}
	heur += sum * resOdds

	// Resources in hand
	/*
		sum = 0
		for _, quant := range player.res {
			sum += float32(quant)
		}
		heur += sum * resInHand
	*/

	// Longest road
	heur += float32(player.longestRoad) / float32(game.longestRoad) * laDistance
	// Largest army
	kin := 0
	for i := 0; i < len(player.cards); i++ {
		if player.cards[i] == "kn" {
			kin++
		}
	}
	heur += (float32(player.knightsPlayed) + .5*float32(kin)) / float32(game.largestArmy)
	// Victory points
	heur += float32(player.victoryPoints*player.victoryPoints) * vp
	return heur
}

func playMove(m Move, pn int, g Game) Game {

	game := fullCopyGame(g)
	player := *game.players[pn]
	move := copyMove(m)

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
	for k, v := range move.gain {
		player.res[k] += v
	}
	for i := 0; i < len(move.newNodes); i++ {
		for j := 0; j < 3; j++ {
			if game.board.nodes[i].edges[j].inUse {
				if game.board.nodes[move.newNodes[i].index].edges[j].nodes[0].index == move.newNodes[i].index {
					move.newNodes[i].edges = game.board.nodes[move.newNodes[i].index].edges[j].nodes[0].edges
					move.newNodes[i].tiles = game.board.nodes[move.newNodes[i].index].edges[j].nodes[0].tiles
					game.board.nodes[move.newNodes[i].index].edges[j].nodes[0] = &move.newNodes[i]
				} else {
					move.newNodes[i].edges = game.board.nodes[move.newNodes[i].index].edges[j].nodes[1].edges
					move.newNodes[i].tiles = game.board.nodes[move.newNodes[i].index].edges[j].nodes[1].tiles
					game.board.nodes[move.newNodes[i].index].edges[j].nodes[1] = &move.newNodes[i]
				}
			}
		}
		game.board.nodes[move.newNodes[i].index] = &move.newNodes[i]
		game.players[pn].villages = append(game.players[pn].villages, game.board.nodes[move.newNodes[i].index])
		player.victoryPoints += move.newNodes[i].building
	}
	// Adds the new roads, as well as searches if there's a newest longest road
	var wg sync.WaitGroup
	var roadLens []*int
	for i := 0; i < len(move.newEdges); i++ {
		move.newEdges[i].nodes[0] = game.board.nodes[move.newEdges[i].nodes[0].index]
		move.newEdges[i].nodes[1] = game.board.nodes[move.newEdges[i].nodes[1].index]
		game.board.edges[move.newEdges[i].index].road = move.newEdges[i].road
		game.players[pn].roads = append(game.players[pn].roads, game.board.edges[move.newEdges[i].index])
		val := 0
		roadLens = append(roadLens, &val)
		wg.Add(1)
		go walkRoad(&wg, game.board.edges[move.newEdges[i].index], player.number, &val, make(map[int]bool))
	}
	wg.Wait()
	if len(roadLens) > 0 {
		greatest := 0
		for i := 1; i < len(roadLens); i++ {
			if *roadLens[i] > *roadLens[greatest] {
				greatest = i
			}
		}
		if *roadLens[greatest] > game.longestRoad {
			if game.lrPlayer != -1 {
				game.players[game.lrPlayer].victoryPoints -= 2
			}
			player.victoryPoints += 2
			game.longestRoad = *roadLens[greatest]
			game.lrPlayer = player.number
		}
		if *roadLens[greatest] > player.longestRoad {
			player.longestRoad = *roadLens[greatest]
		}
	}
	if player.knightsPlayed > game.largestArmy {
		if game.laPlayer != -1 {
			game.players[game.laPlayer].victoryPoints -= 2
		}
		player.victoryPoints += 2
		game.largestArmy = player.knightsPlayed
		game.laPlayer = player.number
	}

	for i := 0; i < move.cardsBought; i++ {
		card := game.cards[0]
		game.cards = game.cards[1:]
		if card == "vp" {
			player.victoryPoints += 1
		} else {
			player.cards = append(player.cards, card)
		}
	}

	game.players[player.number] = &player

	return game
}

func placeStart(board *Board, player *Player) {
	var wg sync.WaitGroup
	heur := make(map[int]*float32)
	for i := 0; i < len(board.nodes); i++ {
		if isBuildable(board.nodes[i]) {
			var f float32
			heur[i] = &f
			wg.Add(1)
			go settleHeur(&wg, board, board.nodes[i], player, &f)
		}
	}

	wg.Wait()
	bestIn := -1
	for i := range heur {
		if bestIn == -1 || *heur[i] > *heur[bestIn] {
			bestIn = i
		}
	}

	board.nodes[bestIn].building = 1
	board.nodes[bestIn].owner = player.number
	player.villages = append(player.villages, board.nodes[bestIn])

	rand.Seed(time.Now().UnixNano())
	for true {
		r := rand.Intn(3)
		if board.nodes[bestIn].edges[r].inUse {
			board.nodes[bestIn].edges[r].road = player.number
			player.roads = append(player.roads, board.nodes[bestIn].edges[r])
			break
		}
	}
	// TODO: If I have time add some actual logic for this
	/*
		if len(player.villages) > 1 {
			if player.villages[0].index > player.villages[1].index {

			} else {

			}
		} else {

		}
	*/
	fmt.Printf("Player %d claimed Node %d \n", player.number, board.nodes[bestIn].index)
}

func isBuildable(node *Node) bool {
	for i := 0; i < 3; i++ {
		if node.edges[i].inUse {
			if node.edges[i].nodes[0].building != 0 {
				return false
			}
			if node.edges[i].nodes[1].building != 0 {
				return false
			}
		}
	}
	return true
}

func settleHeur(wg *sync.WaitGroup, board *Board, node *Node, player *Player, heur *float32) {
	defer wg.Done()
	var count int32
	explored := make(map[int]bool)
	var wg2 sync.WaitGroup
	var mut sync.Mutex
	wg2.Add(1)
	searchVil(&wg2, &mut, node, 4, &count, &explored)
	wg2.Wait()
	res := make(map[string]float32)
	for i := 0; i < len(node.tiles); i++ {
		if _, ok := res[node.tiles[i].res]; !ok {
			res[node.tiles[i].res] = diceOdds[node.tiles[i].roll-1]
		} else {
			res[node.tiles[i].res] += diceOdds[node.tiles[i].roll-1]
		}
	}
	for i := 0; i < len(player.villages); i++ {
		for j := 0; j < len(player.villages[i].tiles); j++ {
			res[player.villages[i].tiles[j].res] += diceOdds[player.villages[i].tiles[j].roll-1]
		}
	}

	var totalOdds float32 = 0
	resCount := 0
	for name, val := range res {
		if name != "D" {
			totalOdds += val * 15
			resCount++
		}
	}

	*heur = (totalOdds) + (float32(resCount) * 2) - float32(count*2)
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

func walkRoad(wg *sync.WaitGroup, edge *Edge, playerNum int, count *int, explored map[int]bool) {
	defer wg.Done()
	if _, ok := explored[edge.index]; !ok {
		return
	}
	c := explored
	c[edge.index] = true
	var counts []*int
	var wg2 sync.WaitGroup

	edge1 := edge.nodes[0].edges[0]
	edge2 := edge.nodes[0].edges[1]
	edge3 := edge.nodes[0].edges[2]
	ok1 := explored[edge1.index]
	ok2 := explored[edge2.index]
	ok3 := explored[edge3.index]
	if boolToInt(ok1)+boolToInt(ok2)+boolToInt(ok3) == 2 {
		if edge1.inUse && edge1.road == playerNum {
			co := 0
			counts = append(counts, &co)
			wg2.Add(1)
			go walkRoad(&wg2, edge1, playerNum, &co, c)
		}
		if edge2.inUse && edge2.road == playerNum {
			co := 0
			counts = append(counts, &co)
			wg2.Add(1)
			go walkRoad(&wg2, edge2, playerNum, &co, c)
		}
		if edge3.inUse && edge3.road == playerNum {
			co := 0
			counts = append(counts, &co)
			wg2.Add(1)
			go walkRoad(&wg2, edge3, playerNum, &co, c)
		}
	}
	edge1 = edge.nodes[1].edges[0]
	edge2 = edge.nodes[1].edges[1]
	edge3 = edge.nodes[1].edges[2]
	ok1 = explored[edge1.index]
	ok2 = explored[edge2.index]
	ok3 = explored[edge3.index]
	if boolToInt(ok1)+boolToInt(ok2)+boolToInt(ok3) == 2 {
		if edge1.inUse && edge1.road == playerNum {
			co := 0
			counts = append(counts, &co)
			wg2.Add(1)
			go walkRoad(&wg2, edge1, playerNum, &co, c)
		}
		if edge2.inUse && edge2.road == playerNum {
			co := 0
			counts = append(counts, &co)
			wg2.Add(1)
			go walkRoad(&wg2, edge2, playerNum, &co, c)
		}
		if edge3.inUse && edge3.road == playerNum {
			co := 0
			counts = append(counts, &co)
			wg2.Add(1)
			go walkRoad(&wg2, edge3, playerNum, &co, c)
		}
	}
	wg2.Wait()
	bestVal := 0
	for i := 0; i < len(counts); i++ {
		if *counts[i] > bestVal {
			bestVal = *counts[i]
		}
	}
	*count = 1 + bestVal
}

// An annoying helper function because bools aren't ints
func boolToInt(val bool) int {
	if val {
		return 1
	} else {
		return 0
	}
}

func copyMove(src Move) Move {
	m := src

	m.cost = make(map[string]int)
	m.gain = make(map[string]int)
	for k, val := range src.cost {
		m.cost[k] = val
	}
	for k, val := range src.gain {
		m.gain[k] = val
	}
	copy(m.newNodes, src.newNodes)
	copy(m.newEdges, src.newEdges)

	return m
}

func copyPlayer(src Player) Player {
	p := src

	p.res = make(map[string]int)
	for k, val := range src.res {
		p.res[k] = val
	}
	copy(p.villages, src.villages)
	copy(p.roads, src.roads)
	copy(p.cards, src.cards)

	return p
}

func fullCopyGame(src Game) Game {
	g := Game{
		longestRoad: src.longestRoad,
		lrPlayer:    src.lrPlayer,
		largestArmy: src.largestArmy,
		laPlayer:    src.largestArmy,
	}
	g.board = copyBoard(src.board)
	copy(g.cards, src.cards)

	for i := 0; i < 4; i++ {
		p := copyPlayer(*src.players[i])
		for j := 0; j < len(p.villages); j++ {
			p.villages[j] = g.board.nodes[p.villages[j].index]
		}
		for j := 0; j < len(p.roads); j++ {
			p.roads[j] = g.board.edges[p.roads[j].index]
		}
		g.players[i] = &p
	}

	return g
}

func copyBoard(src Board) Board {
	b := Board{}

	for i := 0; i < len(src.tiles); i++ {
		c := *src.tiles[i]
		b.tiles = append(b.tiles, &c)
	}
	for i := 0; i < len(src.nodes); i++ {
		c := *src.nodes[i]
		b.nodes = append(b.nodes, &c)
	}
	for i := 0; i < len(src.edges); i++ {
		c := *src.edges[i]
		b.edges = append(b.edges, &c)
	}
	return b
}
