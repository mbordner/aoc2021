package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"aoc2021/common/graph"
	"aoc2021/common/graph/djikstra"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type BoardPosType int

const (
	Wall BoardPosType = iota
	Goal
	Hall
	AboveGoal
)

const (
	PieceTypeCount = 4
	Debug          = false
)

type Pieces map[geom.Pos]*Piece
type PieceType string
type States []*State

var (
	start                   = time.Now()
	pieceTypeCostMultiplier = map[PieceType]int{"A": 1, "B": 10, "C": 100, "D": 1000}
	reNodeChar              = regexp.MustCompile(`#|\.|[A-Z]`)
	memo                    = make(map[string][]States)
)

func (ps Pieces) Clone() Pieces {
	nps := make(Pieces)
	for k, v := range ps {
		nps[k] = v
	}
	return nps
}

func (ss States) Cost() int {
	cost := 0
	for _, s := range ss {
		cost += s.cost
	}
	return cost
}

func (ss States) String() string {
	strs := make([]string, 0, len(ss))
	for _, s := range ss {
		strs = append(strs, s.String())
	}
	return strings.Join(strs, "-----\n") + "\n"
}

func (ps Pieces) CloneWithNewPos(piece *Piece, pos geom.Pos) Pieces {
	nps := make(Pieces)
	for k, v := range ps {
		if v == piece {
			nps[pos] = piece
		} else {
			nps[k] = v
		}
	}
	return nps
}

func (ps Pieces) String() string {
	spots := make(map[PieceType]geom.Positions)
	types := make([]PieceType, 0, 4)
	for pos, piece := range ps {
		if a, e := spots[piece.pieceType]; e {
			spots[piece.pieceType] = append(a, pos)
		} else {
			types = append(types, piece.pieceType)
			spots[piece.pieceType] = []geom.Pos{pos}
		}
	}

	for _, arr := range spots {

		sort.Slice(arr, func(i, j int) bool {
			if arr[i].X == arr[j].X && arr[i].Y == arr[j].Y {
				if arr[i].Z < arr[j].Z {
					return true
				}
			}
			if arr[i].X == arr[j].X {
				if arr[i].Y < arr[j].Y {
					return true
				}
			}
			if arr[i].X < arr[j].X {
				return true
			}
			return false

		})
	}

	sort.Slice(types, func(i, j int) bool {
		return types[i] < types[j]
	})

	strs := make([]string, len(types), len(types))
	for i, t := range types {
		strs[i] = fmt.Sprintf(" [%s: %s] ", t, spots[t].String())
	}

	return strings.Join(strs, ",")
}

type Move struct {
	cost int
	pos  geom.Pos
}

type Piece struct {
	id             int
	costMultiplier int
	pieceType      PieceType
}

func (p *Piece) getMoves(curPos geom.Pos, pieces Pieces, board *Board) []Move {

	if Debug {
		fmt.Println("! calling getMoves for piece ", p.pieceType, curPos, " with pieces: ", pieces)
	}
	curNode := board.g.GetNode(curPos)
	moves := make([]Move, 0, 25)

	curType := curNode.GetProperty("type").(BoardPosType)

	if Debug {
		fmt.Println("> starting to find moves for piece ", p.pieceType, " at position ", curPos)
	}

	pathBlocked := func(nodes []*graph.Node) bool {
		for _, n := range nodes {
			pos := n.GetID().(geom.Pos)
			if _, occupied := pieces[pos]; occupied {
				return true
			}
		}
		return false
	}

	pieceOnGoal := false
	// add all moves for our desired goal whether we are in the hall or in a goal
	// the nodes returned for our goal type will be sorted so the lowest is returned first
	// so we only need to look at the first unoccupied
	var piecesStr string
	if Debug {
		piecesStr = pieces.String()
	}
	for _, n := range board.goals[p.pieceType] {
		nPos := n.GetID().(geom.Pos)
		if Debug {
			fmt.Println("checking for goal target at position: ", nPos)
		}
		if _, occupied := pieces[nPos]; !occupied {
			path, pathCost := board.shortestPaths[curNode].GetShortestPath(n)
			if !pathBlocked(path) {
				if Debug {
					fmt.Println("found goal target at position: ", nPos)
				}
				moves = append(moves, Move{cost: p.costMultiplier * int(pathCost), pos: nPos})
				break // if we can go to the lowest goal, break and lock out going to an above goal spot
			}
		} else {
			if pieces[nPos].id == p.id {
				if Debug {
					fmt.Println("position ", curPos, " is occupied by this piece, so marking it as goal")
				}
				pieceOnGoal = true
				break // the occupied space is us.. so we are at the goal
			}
			if pieces[nPos].pieceType != p.pieceType {
				if Debug {
					fmt.Println("this position is occupied by another type (", pieces[nPos].pieceType, ") so goals are not an option")
				}
				break // this break will only be necessary when we're looking at the lower goals
				// if a lower one is occupied by a piece of another type, we can't go to a spot above
			}
		}
	}
	if Debug {
		checkStr := pieces.String()
		if checkStr != piecesStr {
			fmt.Println("pieces changed.. not supposed to happen")
		}
	}

	// if we're in the hall, we can only go to goal, however if we are in a goal, we can also go to the hall spots
	if curType != Hall && !pieceOnGoal {
		if Debug {
			fmt.Println("piece is not on a hall, so checking for available hall spots")
		}
		for _, n := range board.nodeTypes[Hall] {
			nPos := n.GetID().(geom.Pos)
			if _, occupied := pieces[nPos]; !occupied {
				path, pathCost := board.shortestPaths[curNode].GetShortestPath(n)
				if !pathBlocked(path) {
					if Debug {
						fmt.Println("found a hall target at: ", nPos)
					}
					moves = append(moves, Move{cost: p.costMultiplier * int(pathCost), pos: nPos})
				}
			}
		}
	}
	if Debug {
		checkStr := pieces.String()
		if checkStr != piecesStr {
			fmt.Println("pieces changed.. not supposed to happen")
		}
	}

	if Debug {
		fmt.Println("< returning after ", len(moves), " moves found")
	}

	return moves

}

type Board struct {
	g             *graph.Graph
	startPieces   Pieces
	nextID        int
	bb            geom.BoundingBox
	shortestPaths map[*graph.Node]djikstra.ShortestPaths
	nodeTypes     map[BoardPosType][]*graph.Node
	goals         map[PieceType][]*graph.Node
}

func (b *Board) nextPieceID() int {
	id := b.nextID
	b.nextID++
	return id
}

func (b *Board) addPiece(pieceType PieceType, pos geom.Pos) {
	p := new(Piece)
	p.id = b.nextPieceID()
	p.pieceType = pieceType
	p.costMultiplier = pieceTypeCostMultiplier[pieceType]
	b.startPieces[pos] = p
}

type State struct {
	cost   int
	pieces Pieces
	board  *Board
}

func (s State) String() string {
	height := s.board.bb.MaxY + 1
	width := s.board.bb.MaxX + 1

	lines := make([]string, 0, height)

	for y := 0; y < height; y++ {
		line := make([]byte, width, width)

		for x := 0; x <= s.board.bb.MaxX; x++ {
			pos := geom.Pos{Y: y, X: x}
			if piece, occupied := s.pieces[pos]; occupied {
				line[x] = piece.pieceType[0]
			} else {
				n := s.board.g.GetNode(pos)
				if n != nil {
					posType := n.GetProperty("type").(BoardPosType)
					if posType == Wall {
						line[x] = '#'
					} else if posType == Goal || posType == AboveGoal || posType == Hall {
						line[x] = '.'
					}
				} else {
					line[x] = ' '
				}

			}
		}

		lines = append(lines, string(line))
	}

	return strings.Join(lines, "\n") + fmt.Sprintf("cost: %d\n", s.cost)
}

func newState(board *Board, cost int, pieces Pieces) *State {
	s := new(State)
	s.board = board
	s.cost = cost
	s.pieces = pieces
	return s
}

func (p *Piece) isGoal(board *Board, state *State, pos geom.Pos) bool {
	for _, g := range board.goals[p.pieceType] {
		gPos := g.GetID().(geom.Pos)
		if gp, occupied := state.pieces[gPos]; occupied {
			if gp.pieceType != p.pieceType {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (s *State) isGoal() bool {
	for pos, piece := range s.pieces {
		if !piece.isGoal(s.board, s, pos) {
			return false
		}
	}
	return true
}

func (s *State) getStatesToGoal() []States {
	if Debug {
		fmt.Println("starting a new state :", s.pieces.String())
		fmt.Println(s)
	}

	memoKey := fmt.Sprintf("%d|%s", s.cost, s.pieces.String())
	if states, exist := memo[memoKey]; exist {
		if Debug {
			fmt.Println("found a memo string, returning")
		}
		return states
	}

	if s.isGoal() {
		memo[memoKey] = []States{{s}}
		return memo[memoKey]
	}

	statesToGoal := make([]States, 0, 100)

	for pos, piece := range s.pieces {
		if !piece.isGoal(s.board, s, pos) {
			moves := piece.getMoves(pos, s.pieces, s.board)
			for _, move := range moves {
				statesToGoalFromPiece := newState(s.board, move.cost, s.pieces.CloneWithNewPos(piece, move.pos)).getStatesToGoal()
				if statesToGoalFromPiece != nil {
					for _, states := range statesToGoalFromPiece {
						statesToGoal = append(statesToGoal, append(States{s}, states...))
					}
				}
			}
		}
	}

	if len(statesToGoal) > 0 {

		var minStates States
		minCost := statesToGoal[0].Cost()
		minStates = statesToGoal[0]

		for _, states := range statesToGoal {
			cost := states.Cost()
			if cost < minCost {
				minCost = cost
				minStates = states
			}
		}

		statesToGoal = []States{minStates}

		memo[memoKey] = statesToGoal

		if Debug {
			fmt.Sprintf("<< found some paths to goal, but exhausted our piece checks.. popping back")
		}
	} else {

		if Debug {
			fmt.Sprintf("<< didn't find a goal from this state, popping back")
		}
		memo[memoKey] = nil

	}

	return memo[memoKey]
}

func getInitialBoard(filename string) *Board {

	b := new(Board)
	b.g = graph.NewGraph()
	b.startPieces = make(Pieces)
	b.shortestPaths = make(map[*graph.Node]djikstra.ShortestPaths)
	b.nodeTypes = make(map[BoardPosType][]*graph.Node)
	b.goals = make(map[PieceType][]*graph.Node)

	for k := range pieceTypeCostMultiplier {
		b.goals[PieceType(k)] = make([]*graph.Node, 0, PieceTypeCount)
	}

	lines, _ := file.GetLines(filename)

	pieceXs := make(map[int]PieceType)
	pieceCol := 'A' - 1

	for y, line := range lines {

		for x, char := range line {

			if reNodeChar.MatchString(string(char)) {
				pos := geom.Pos{X: x, Y: y}
				b.bb.Extend(pos)
				node := b.g.CreateNode(pos)

				if char == '#' {
					node.SetTraversable(false)
					node.AddProperty("type", Wall)
				} else if char == '.' {
					node.AddProperty("type", Hall)
				} else {
					pieceType := PieceType(char)
					b.addPiece(pieceType, pos)

					node.AddProperty("type", Goal)

					var goalType PieceType
					if gt, set := pieceXs[x]; set {
						goalType = gt
					} else {
						pieceCol++
						goalType = PieceType(pieceCol)
						pieceXs[x] = goalType
					}

					node.AddProperty("goal", goalType)
				}
			}

		}
	}

	// mark off above goal nodes

	for x, pieceType := range pieceXs {
		y := 0
		n := b.g.GetNode(geom.Pos{X: x, Y: y})
		for n != nil {
			if pt := n.GetProperty("type").(BoardPosType); pt != Wall && pt != Goal {
				n.AddProperty("type", AboveGoal)
				n.AddProperty("goal", pieceType)
				break
			}
			y++
			n = b.g.GetNode(geom.Pos{X: x, Y: y})
		}
	}

	// add edges

	for _, n := range b.g.GetTraversableNodes() {

		pos := n.GetID().(geom.Pos)

		var o *graph.Node
		var e *graph.Edge

		// left
		o = b.g.GetNode(pos.Transform(-1, 0, 0))
		if o != nil && o.IsTraversable() {
			e = n.AddEdge(o, float64(1))
			e.AddProperty("dir", geom.West)
		}

		// right
		o = b.g.GetNode(pos.Transform(1, 0, 0))
		if o != nil && o.IsTraversable() {
			e = n.AddEdge(o, float64(1))
			e.AddProperty("dir", geom.East)
		}

		// above
		o = b.g.GetNode(pos.Transform(0, -1, 0))
		if o != nil && o.IsTraversable() {
			e = n.AddEdge(o, float64(1))
			e.AddProperty("dir", geom.North)
		}

		// below
		o = b.g.GetNode(pos.Transform(0, 1, 0))
		if o != nil && o.IsTraversable() {
			e = n.AddEdge(o, float64(1))
			e.AddProperty("dir", geom.South)
		}
	}

	// generate stats

	for _, n := range b.g.GetTraversableNodes() {

		nodeType := n.GetProperty("type").(BoardPosType)
		if _, e := b.nodeTypes[nodeType]; !e {
			b.nodeTypes[nodeType] = []*graph.Node{}
		}
		b.nodeTypes[nodeType] = append(b.nodeTypes[nodeType], n)

		if nodeType == Goal {
			pieceType := n.GetProperty("goal").(PieceType)
			if _, e := b.goals[pieceType]; !e {
				b.goals[pieceType] = []*graph.Node{}
			}
			b.goals[pieceType] = append(b.goals[pieceType], n)
		}

		b.shortestPaths[n] = djikstra.GenerateShortestPaths(b.g, n)
	}

	for _, nodes := range b.goals {
		sort.Slice(nodes, func(i, j int) bool {
			iPos := nodes[i].GetID().(geom.Pos)
			jPos := nodes[j].GetID().(geom.Pos)
			if iPos.Y > jPos.Y {
				return true
			}
			return false
		})
	}

	return b

}

func main() {
	board := getInitialBoard("../data2.txt")

	initialState := newState(board, 0, board.startPieces.Clone())
	statesToGoal := initialState.getStatesToGoal()

	if statesToGoal == nil {
		fmt.Println("no path to goal found")
	} else {
		var minStates States
		minCost := statesToGoal[0].Cost()
		minStates = statesToGoal[0]

		for _, states := range statesToGoal {
			cost := states.Cost()
			if cost < minCost {
				minCost = cost
				minStates = states
			}
		}

		fmt.Println(minStates, fmt.Sprintf("total cost: %d", minCost))
	}

	duration := time.Since(start)
	fmt.Println(duration)

}
