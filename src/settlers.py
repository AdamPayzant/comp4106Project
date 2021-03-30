from structures import Map
from players import Player, Human, AI

import random

class Game:
    players = []
    map
    def __init__(self, human=False):
        self.map = Map()
        if not human:
            self.players.append(AI(0))
            self.players.append(AI(1))
            self.players.append(AI(2))
            self.players.append(AI(3))
        else:
            self.players.append(AI(0))
            self.players.append(AI(1))
            self.players.append(AI(2))
            self.players.append(Human(3))

    def start(self) -> int:
        i = random.randint(0,3)
        cur = self.players[i]
        self.__place__(i)

        while cur.victoryPoints < 10:
            cur = self.players[i]
            roll = random.randint(1,6) + random.randint(1,6)
            if roll == 7:
                # Do bandit handling
                pass
            else:
                for player in self.players:
                    player.distRoll(roll)

            cur.play()

            i += 1
            if i > 3:
                i = 0

        return cur.index

    def __place__(self, i):
        pass