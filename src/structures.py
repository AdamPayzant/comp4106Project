import random
# just a list of all tile types
TILES = [
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
    "de"
]
# List of all possible rolls
ROLLS = [
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
    12
]
# List of each row length
ROWLENS = [
    7,
    9,
    11,
    11,
    9,
    7
]
# An offset for tile calculations
TILEOFF = [
    0,
    3,
    7,
    7,
    12,
    16
]

class Tile:
    roll = 0 # The roll to get this resource
    res = "" # The resource for this tile
    bandit = False # Does the tile have a bandit

    def __init__(self, res: str, roll: int):
        self.res = res
        self.roll = roll

class Port:
    res = "?"
    cost = 0
    def __init__(self, res, cost):
        self.res = res
        self.cost = cost

class Vertex:
    nodes = []
    road = 0 # player numbers, 0 means no road
    index = 0

class Node:
    index = 0
    port = None
    building = 0 # 0 for nothing, 1 for settlement, 2 for city
    vertex = []
    tiles = []
    def __init__(self, index, tiles, port = None):
        self.index = index
        self.port = port
        self.tiles = tiles

class Map:
    tiles = []
    nodes = []

    def __init__(self):
        self.__genMap()

    def __genMap(self):
        # Generate the Tiles
        tilecp = TILES
        rollcp = ROLLS
        while tilecp:
            random.shuffle(tilecp)
            random.shuffle(rollcp)
            t = tilecp.pop(0)
            if t != "de":
                r = rollcp.pop(0)
            else:
                r = 7
            self.tiles.append(Tile(t, r))

        nodes = []
        # I know, I know, this is really bad
        # I was trying to grind this part out quickly and it ballooned pretty bad
        for row in range(0, len(ROWLENS)):
            nodes.append([])
            # The first row
            if row == 0:
                for j in range(0, ROWLENS[row]):
                    print(j)
                    t = []
                    if j < ROWLENS[row]-1:
                        t.append(self.tiles[j//2 + TILEOFF[row]])
                    if j > 0 and row % 2 == 0:
                        t.append(self.tiles[j//2 - 1 +TILEOFF[row]])
                    sum = 0
                    for i in range(0,row):
                        sum += ROWLENS[i]
                    node = Node(sum + j, t)

                    if j == 0:
                        vert1 = None
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    elif j == ROWLENS[row] - 1:
                        vert1 = nodes[row][j - 1].vertex[2]
                        vert2 = Vertex
                        vert3 = None

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                    elif j % 2 == 0:
                        vert1 = nodes[row][j - 1].vertex[2]
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert3.nodes.append(node)
                    else:
                        vert1 = nodes[row][j - 1].vertex[2]
                        vert2 = None
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert3.nodes.append(node)

                    node.vertex = [vert1, vert2, vert3]
                    nodes[row].append(node)
            # The last row
            elif row == len(ROWLENS) - 1:
                for j in range(0, ROWLENS[row]):
                    t = []
                    if j < ROWLENS[row]-1:
                        t.append(self.tiles[j//2 + TILEOFF[row]])
                    if j > 0 and row % 2 == 0:
                        t.append(self.tiles[j//2 - 1 +TILEOFF[row]])
                    sum = 0
                    for i in range(0,row):
                        sum += ROWLENS[i]
                    node = Node(sum + j, t)

                    if j == 0:
                        vert1 = None
                        vert2 = nodes[row-1][j-1].vertex[1]
                        vert3 = Vertex()

                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    elif j == ROWLENS[row] - 1:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = nodes[row-1][j-1].vertex[1]
                        vert3 = None

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                    elif j % 2 == 0:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = None
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert3.nodes.append(node)
                    else:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)

                    node.vertex = [vert1, vert2, vert3]
                    nodes[row].append(node)
            # Row Length is increasing
            elif ROWLENS[row] > ROWLENS[row - 1]:
                for j in range(0, ROWLENS[row]):
                    t = []
                    if j == 0:
                        t.append(self.tiles[TILEOFF[row]])
                    elif j == ROWLENS[row] - 1:
                        t.append(self.tiles[j//2 + TILEOFF[row] - 1])
                    elif j % 2 == 0:
                        t.append(self.tiles[j//2 + TILEOFF[row] - 1])
                        t.append(self.tiles[j//2 + TILEOFF[row]])
                        t.append(self.tiles[j//2 + TILEOFF[row - 1] - 1])
                    else:
                        if j == 1:
                            t.append(self.tiles[TILEOFF[row]])
                            t.append(self.tiles[TILEOFF[row-1]])
                        elif j == ROWLENS[row] - 2:
                            t.append(self.tiles[j//2 + TILEOFF[row]])
                            t.append(self.tiles[TILEOFF[row] - 1])
                        else:
                            t.append(self.tiles[j//2 + TILEOFF[row-1] -1])
                            t.append(self.tiles[j//2 + TILEOFF[row-1]])
                            t.append(self.tiles[j//2 + TILEOFF[row]])
                    sum = 0
                    for i in range(0,row):
                        sum += ROWLENS[i]
                    node = Node(sum + j, t)

                    if j == 0:
                        vert1 = None
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    elif j % 2 == 0:
                        vert1 = nodes[row][j - 1].vertex[2]
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    else:
                        vert1 = nodes[row][j - 1].vertex[2]
                        vert2 = nodes[row-1][j-1].vertex[1]
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)

                    node.vertex = [vert1, vert2, vert3]
                    nodes[row].append(node)
            # Row length is staying the same
            elif ROWLENS[row] == ROWLENS[row - 1]:
                for j in range(0, ROWLENS[row]):
                    t = []
                    if j == 0:
                        t.append(self.tiles[TILEOFF[row]])
                    elif j == ROWLENS[row] - 1:
                        t.append(self.tiles[TILEOFF[row+1] - 1])
                    elif j % 2 == 0:
                        t.append(self.tiles[j//2 + TILEOFF[row]])
                        t.append(self.tiles[j//2 + TILEOFF[row] - 1])
                        t.append(self.tiles[j//2 + TILEOFF[row+1] - 1])
                    else:
                        if j == 1:
                            t.append(self.tiles[TILEOFF[row+1]])
                            t.append(self.tiles[TILEOFF[row]])
                        elif j == ROWLENS[row] - 2:
                            t.append(self.tiles[j//2 - 1 + TILEOFF[row+1]])
                            t.append(self.tiles[j//2 + TILEOFF[row]])
                        else:
                            t.append(self.tiles[j//2 + TILEOFF[row] - 1])
                            t.append(self.tiles[j//2 + TILEOFF[row+1] -1])
                            t.append(self.tiles[j//2 + TILEOFF[row+1]])
                    sum = 0
                    for i in range(0,row):
                        sum += ROWLENS[i]
                    node = Node(sum + j, t)

                    if j == 0:
                        vert1 = None
                        vert2 = nodes[row-1][j].vertex[1]
                        vert3 = Vertex()

                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    elif j % 2 == 0:
                        vert1 = nodes[row][j - 1].vertex[2]
                        vert2 = nodes[row-1][j].vertex[1]
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    else:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)

                    node.vertex = [vert1, vert2, vert3]
                    nodes[row].append(node)
            # Row length is decreasing
            else:
                for j in range(0, ROWLENS[row]):
                    t = []
                    if j == 0:
                        t.append(self.tiles[TILEOFF[row]])
                    elif j == ROWLENS[row] - 1:
                        t.append(self.tiles[TILEOFF[row+1] - 1])
                    elif j % 2 == 0:
                        t.append(self.tiles[j//2 + TILEOFF[row]])
                        t.append(self.tiles[j//2 + TILEOFF[row] - 1])
                        t.append(self.tiles[j//2 + TILEOFF[row+1] - 1])
                    else:
                        if j == 1:
                            t.append(self.tiles[TILEOFF[row+1]])
                            t.append(self.tiles[TILEOFF[row]])
                        elif j == ROWLENS[row] - 2:
                            t.append(self.tiles[j//2 - 1 + TILEOFF[row+1]])
                            t.append(self.tiles[j//2 + TILEOFF[row]])
                        else:
                            t.append(self.tiles[j//2 + TILEOFF[row] - 1])
                            t.append(self.tiles[j//2 + TILEOFF[row+1] -1])
                            t.append(self.tiles[j//2 + TILEOFF[row+1]])
                    sum = 0
                    for i in range(0,row):
                        sum += ROWLENS[i]
                    node = Node(sum + j, t)

                    if j == 0:
                        vert1 = None
                        vert2 = nodes[row-1][j-1].vertex[1]
                        vert3 = Vertex()

                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    elif j == ROWLENS[row] - 1:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = nodes[row-1][j-1].vertex[1]
                        vert3 = None

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                    elif j % 2 == 0:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = nodes[row-1][j-1].vertex[1]
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)
                    else:
                        vert1 = nodes[row][j-1].vertex[2]
                        vert2 = Vertex()
                        vert3 = Vertex()

                        vert1.nodes.append(node)
                        vert2.nodes.append(node)
                        vert3.nodes.append(node)

                    node.vertex = [vert1, vert2, vert3]
                    nodes[row].append(node)
        self.nodes = nodes

    def printMap(self):
        print(self.nodes)