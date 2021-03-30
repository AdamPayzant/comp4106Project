from structures import Map

import random
import copy
import threading

# The maximum depth we minMax to
# NOTE: 1 depth = 1 full pass of the table
MAXDEPTH = 3

class Player:
    victoryPoints = 0
    index = 0

    def distRoll(self, roll):
        pass
    def play(self, map: Map):
        pass
    def predict(self, map: Map, hostID, depth, val):
        pass
    def heuristic(self, map: Map, player: Player):
        pass

class Human(Player):
    victoryPoints = 0
    index = 0
    others = []
    cards = []
    res = []

    villages = []
    knights = []
    longest_road = False

    def __init__(self, index):
        self.index = index

    def distRoll(self, roll):
        for vil in self.villages:
            for tile in vil.tiles:
                if tile.roll == roll and not tile.bandit:
                    self.res.extend([tile.res for i in range(0, vil.building)])

    def play(self, map: Map):
        map.printMap()
        # TODO: Give the player options to play

    def continueMove(self, map, move):
        pass

    def predict(self, map: Map, hostID, depth, val):
        if depth >= MAXDEPTH:
            return 0

    def heuristic(self, map: Map, player: Player):
        pass

    # Removes all of a specified resource and returns the count
    def monopoly(self, res) -> int:
        count = self.res.count(res)
        self.res = [i for i in self.res if i != res]
        return count
        

class AI(Player):
    victoryPoints = 0
    index = 0
    others = []
    res = []
    cards = []

    villages = []
    knights = []
    longest_road = False

    def __init__(self, index):
        self.index = index

    def distRoll(self, roll):
        for vil in self.villages:
            for tile in vil.tiles:
                if tile.roll == roll and not tile.bandit:
                    self.res.extend([tile.res for i in range(0, vil.building)])

    def setOthers(self, others):
        self.others = others

    def play(self, map: Map):
        map.printMap()

        # So playing and predicting lends itself really well to multithreading. 
        # Basically I pop off threads for every potential line
        # TODO: If memory usage gets too extreme, make custom copy method
        threads = []
        moves = []
        # Road case
        if "cl" in self.res and "wo" in self.res:
            roadSite = []
            # TODO: Walk the edges and find possible build sites

        # Settlement case
        if "cl" in self.res and "wh" in self.res and "sh" in self.res and "wo" in self.res:
            buildSite = []
            # TODO: Walk the edges and find possible build sites

        # City case
        if self.res.count("wh") >= 2 and self.res.count("ir") >= 3:
            for vil in self.villages:
                c = vil.softCopy()
                c.building += 1
                m = Move()
                m.nodes.append(c)
                m.cost.extend("wh", "wh", "ir", "ir", "ir")
                threads.append(threading.Thread(
                    target=self.continueMove,
                    args=(map, m)
                ))
                threads[-1].start()

        # Card case
        if "sh" in self.res and "wh" in self.res and "ir" in self.res:
            m = Move()
            m.cost.extend("sh", "wh", "ir")
            m.cardsBought += 1
            threads.append(threading.Thread(
                target=self.continueMove,
                args=(map, m)
            ))
            threads[-1].start()

        # Play card
        if self.cards:
            for card in self.cards:
                m = Move()
                m.card = card
                threads.append(threading.Thread(
                    target=self.continueMove,
                    args=(map, m)
                ))
                threads[-1].start()

        # Pass
        moves.append(Move())
        threads.append(threading.Thread(
            target=copy.deepcopy(
                self.others[self.index+1] 
                if self.index < 3 
                else self.others[0]).predict, 
            args=(map, self.index, 0, moves[-1])))
        threads[-1].start()

        best = [0,0.0]
        for i in range(len(threads)):
            threads[i].join()
            if moves[i].heur > best[1]:
                best = [i, moves[i].heur]
        moves[i].perform(map)

    def continueMove(self, map, move):
        pass

    def predict(self, map: Map, hostID, depth, val):
        if depth >= MAXDEPTH:
            return 0

    def heuristic(self, map: Map, player: Player):
        pass

    # Removes all of a specified resource and returns the count
    def monopoly(self, res) -> int:
        count = self.res.count(res)
        self.res = [i for i in self.res if i != res]
        return count

# Basically a simple way of tracking a move
# Reduces the number of copies needed
class Move:
    cost = []
    edges = []
    nodes = []
    card = None
    cardsBought = 0
    heur = 0.0

    def perform(self, Map):
        pass