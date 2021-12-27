package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"fmt"
)

type SeaCucumberFlockType = geom.Direction
type SeaCucumbers []*SeaCucumber
type SeafloorSpaces []*SeafloorSpace

const (
	EastFacingFlock  = geom.East
	SouthFacingFlock = geom.South
	Debug            = false
)

var (
	FlockTypeMap = map[string]SeaCucumberFlockType{"v": SouthFacingFlock, ">": EastFacingFlock}
)

func (scs SeaCucumbers) GetNextMoves() (SeaCucumbers, SeafloorSpaces) {
	ocs := make(SeaCucumbers, 0, len(scs))
	sss := make(SeafloorSpaces, 0, len(scs))
	for _, sc := range scs {
		ss := sc.GetNextMove()
		if ss != nil {
			sss = append(sss, ss)
			ocs = append(ocs, sc)
		}
	}
	return ocs, sss
}

type SeafloorGrid struct {
	bb        geom.BoundingBox
	grid      map[geom.Pos]*SeafloorSpace
	occupants map[SeaCucumberFlockType]SeaCucumbers
}

func (sg *SeafloorGrid) Print() {
	reverseMap := make(map[SeaCucumberFlockType]string)
	for k, v := range FlockTypeMap {
		reverseMap[v] = k
	}
	maxX := sg.bb.MaxX + 1
	maxY := sg.bb.MaxY + 1
	for y := 0; y < maxY; y++ {
		bytes := make([]byte, maxX, maxX)
		for x := 0; x < maxX; x++ {
			space := sg.GetSpace(geom.Pos{Y: y, X: x})
			if space.IsOccupied() {
				ft := reverseMap[space.GetOccupant().GetFlockType()]
				bytes[x] = ft[0]
			} else {
				bytes[x] = '.'
			}
		}
		fmt.Println(string(bytes))
	}
}

func NewSeafloorGrid() *SeafloorGrid {
	sg := new(SeafloorGrid)
	sg.grid = make(map[geom.Pos]*SeafloorSpace)
	sg.occupants = make(map[SeaCucumberFlockType]SeaCucumbers)
	return sg
}

func (sg *SeafloorGrid) AddSpace(pos geom.Pos, occupant *SeaCucumber) {
	ss := NewSeafloorSpace(sg, pos)
	sg.bb.Extend(pos)
	sg.grid[pos] = ss
	if occupant != nil {
		ss.Enter(occupant)
		ft := occupant.GetFlockType()
		if ocs, exists := sg.occupants[ft]; exists {
			sg.occupants[ft] = append(ocs, occupant)
		} else {
			sg.occupants[ft] = SeaCucumbers{occupant}
		}
	}
}

func (sg *SeafloorGrid) GetOccupants(flockType SeaCucumberFlockType) SeaCucumbers {
	if ocs, exists := sg.occupants[flockType]; exists {
		return ocs
	}
	return nil
}

func (sg *SeafloorGrid) GetSpace(pos geom.Pos) *SeafloorSpace {
	if pos.X > sg.bb.MaxX {
		pos.X = 0
	}
	if pos.Y > sg.bb.MaxY {
		pos.Y = 0
	}
	return sg.grid[pos]
}

type SeafloorSpace struct {
	grid     *SeafloorGrid
	pos      geom.Pos
	occupant *SeaCucumber
}

func NewSeafloorSpace(grid *SeafloorGrid, pos geom.Pos) *SeafloorSpace {
	ss := new(SeafloorSpace)
	ss.grid = grid
	ss.pos = pos
	return ss
}

func (ss *SeafloorSpace) IsOccupied() bool {
	return ss.occupant != nil
}

func (ss *SeafloorSpace) GetOccupant() *SeaCucumber {
	return ss.occupant
}

func (ss *SeafloorSpace) Leave(occupant *SeaCucumber) bool {
	if ss.occupant == occupant {
		ss.occupant = nil
		return true
	}
	return false
}

func (ss *SeafloorSpace) Enter(occupant *SeaCucumber) bool {
	if ss.occupant == nil {
		ss.occupant = occupant
		occupant.Enter(ss)
		return true
	}
	return false
}

func (ss *SeafloorSpace) NextSpace(dir geom.Direction) *SeafloorSpace {
	pos := ss.pos
	switch dir {
	case geom.East:
		pos = pos.Transform(1, 0, 0)
	case geom.South:
		pos = pos.Transform(0, 1, 0)
	}
	return ss.grid.GetSpace(pos)
}

type SeaCucumber struct {
	flockType SeaCucumberFlockType
	location  *SeafloorSpace
}

func NewSeaCucumber(flockType SeaCucumberFlockType) *SeaCucumber {
	sc := new(SeaCucumber)
	sc.flockType = flockType
	return sc
}

func (sc *SeaCucumber) GetFlockType() SeaCucumberFlockType {
	return sc.flockType
}

func (sc *SeaCucumber) GetNextMove() *SeafloorSpace {
	nextPos := sc.location.NextSpace(sc.flockType)
	if !nextPos.IsOccupied() {
		return nextPos
	}
	return nil
}

func (sc *SeaCucumber) Leave() {
	if sc.location != nil {
		sc.location.Leave(sc)
		sc.location = nil
	}
}

func (sc *SeaCucumber) Enter(space *SeafloorSpace) {
	if sc.location == nil {
		sc.location = space
	}
}

func main() {
	sg := getSeafloorGrid("../data.txt")

	numSteps := 0

	for {
		movesMade := 0
		for _, scs := range []SeaCucumbers{sg.GetOccupants(EastFacingFlock), sg.GetOccupants(SouthFacingFlock)} {
			occupants, spaces := scs.GetNextMoves()
			for i := range spaces {
				occupants[i].Leave()
			}
			for i, s := range spaces {
				s.Enter(occupants[i])
			}
			movesMade += len(spaces)
		}

		numSteps++

		if Debug {
			fmt.Println("----")
			sg.Print()
			fmt.Println("----")
		}

		if movesMade == 0 {
			break
		}
	}

	fmt.Println("steps count: ", numSteps)

}

func getSeafloorGrid(filename string) *SeafloorGrid {
	sg := NewSeafloorGrid()

	lines, _ := file.GetLines(filename)

	/*
			str := `..........
		.>v....v..
		.......>..
		..........`

			lines = strings.Split(str,"\n")
	*/

	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			pos := geom.Pos{Y: y, X: x}
			var occupant *SeaCucumber
			char := string(lines[y][x])
			if ft, known := FlockTypeMap[char]; known {
				occupant = NewSeaCucumber(ft)
			}
			sg.AddSpace(pos, occupant)
		}
	}

	return sg
}
