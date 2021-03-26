from structures import *

import random
import threading

# The maximum depth we minMax to
# NOTE: 1 depth = 1 full pass of the table
MAXDEPTH = 3

class Player:
    victoryPoints = 0
    index = 0

    def play(self, map: Map):
        pass
    def predict(self, map: Map) -> float:
        pass

class Human(Player):
    victoryPoints = 0
    index = 0
    others = []
    cards = {}

    def __init__(self, index):
        self.index = index

    def play(self, map: Map):
        map.printMap()
        roll = random.randint(1,6) + random.randint(1,6)
        # TODO: Give the player options to play

    def predict(self, map: Map, i, depth) -> float:
        if depth >= MAXDEPTH:
            return 0
        

class AI(Player):
    victoryPoints = 0
    index = 0
    others = []
    cards = {}

    def __init__(self, index):
        self.index = index

    def setOthers(self, others):
        self.others = others

    def play(self, map: Map):
        map.printMap()

    def predict(self, map: Map, i, depth) -> float:
        if depth >= MAXDEPTH:
            return 0