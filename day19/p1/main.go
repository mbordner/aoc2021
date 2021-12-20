package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	numTransforms = 48
)

var (
	reScanner = regexp.MustCompile(`^---\s+scanner\s+(\d+)\s+---$`)
	signs     = [][]int{{1, 1, 1}, {1, -1, 1}, {-1, 1, 1}, {-1, -1, 1}, {1, 1, -1}, {1, -1, -1}, {-1, 1, -1}, {-1, -1, -1}}
	perms     = [][]int{{0, 1, 2}, {1, 0, 2}, {2, 1, 0}, {1, 2, 0}, {0, 2, 1}, {2, 0, 1}}
)

type offsetData struct {
	offset    Vector
	transform int
}

type Vector struct {
	x int
	y int
	z int
}

func (v Vector) clone() Vector {
	return Vector{x: v.x, y: v.y, z: v.z}
}

func (v Vector) String() string {
	return fmt.Sprintf("{%d,%d,%d}", v.x, v.y, v.z)
}

func abs(v int) int {
	if v < 0 {
		return v * -1
	}
	return v
}

func (v Vector) manhattanDistance(o Vector) int {
	return abs(v.x-o.x) + abs(v.y-o.y) + abs(v.z-o.z)
}

func (v Vector) sub(o Vector) Vector {
	return Vector{x: v.x - o.x, y: v.y - o.y, z: v.z - o.z}
}

func (v Vector) add(o Vector) Vector {
	return Vector{x: v.x + o.x, y: v.y + o.y, z: v.z + o.z}
}

func (v Vector) transform(transformID int) Vector {
	t := v.clone()
	tvs := []int{t.x, t.y, t.z}
	p := transformID % len(perms)
	s := transformID / len(perms)
	t.x = tvs[perms[p][0]] * signs[s][0]
	t.y = tvs[perms[p][1]] * signs[s][1]
	t.z = tvs[perms[p][2]] * signs[s][2]
	return t
}

type Vectors []Vector
type VectorDistances map[Vector]Vector
type MergeVectors map[Vector]int

func (mv *MergeVectors) merge(vs Vectors, offset Vector) {
	for _, v := range vs {
		v = offset.add(v)
		if c, exists := (*mv)[v]; exists {
			(*mv)[v] = c + 1
		} else {
			(*mv)[v] = 1
		}
	}
}

func (mv *MergeVectors) vectors() Vectors {
	vs := make(Vectors, 0, len(*mv))
	for v := range *mv {
		vs = append(vs, v)
	}
	return vs
}

func (vs Vectors) transform(tranformID int) Vectors {
	tvs := make(Vectors, len(vs), len(vs))
	for i := 0; i < len(vs); i++ {
		tvs[i] = vs[i].transform(tranformID)
	}
	return tvs
}

func (vs Vectors) distances() map[Vector]VectorDistances {
	vsDistances := make(map[Vector]VectorDistances)

	for i := 0; i < len(vs); i++ {
		vDistances := make(VectorDistances)
		for j := 0; j < len(vs); j++ {
			if i != j {
				dis := vs[j].sub(vs[i])
				if _, exists := vDistances[dis]; exists {
					panic("already had this distance")
				}
				vDistances[dis] = vs[j]
			}
		}
		vsDistances[vs[i]] = vDistances
	}

	return vsDistances
}

func newVector(l string) Vector {
	v := Vector{}

	tokens := strings.Split(l, ",")
	values := make([]int, len(tokens), len(tokens))
	for i := 0; i < len(tokens); i++ {
		val, _ := strconv.ParseInt(tokens[i], 10, 32)
		values[i] = int(val)
	}

	v.x = values[0]
	v.y = values[1]
	v.z = values[2]

	return v
}

type Scanner struct {
	id           int
	offsetFrom   map[int]offsetData
	vectorReport Vectors
	merged       MergeVectors
}

func newScanner(l string) *Scanner {
	matches := reScanner.FindStringSubmatch(l)
	id, _ := strconv.ParseInt(matches[1], 10, 32)
	s := new(Scanner)
	s.id = int(id)
	s.offsetFrom = make(map[int]offsetData)
	s.vectorReport = make(Vectors, 0, 100)
	s.merged = make(MergeVectors)
	return s
}

func (s *Scanner) merge(vs Vectors, offset Vector) {
	s.merged.merge(vs, offset)
}

func (s *Scanner) add(v Vector) {
	s.vectorReport = append(s.vectorReport, v)
}

func checkOverlap(iDistances, jDistances map[Vector]VectorDistances) (*Vector, VectorDistances, *Vector, VectorDistances) {
	var iVector, jVector *Vector
	var iCommonDistances, jCommonDistances VectorDistances
	maxCommon := 0
outer:
	for vi, vids := range iDistances {
		for vj, vjds := range jDistances {

			commonIDistances := make(VectorDistances)
			commonJDistances := make(VectorDistances)
			for di := range vids {
				if _, exists := vjds[di]; exists {
					commonIDistances[di] = vids[di]
					commonJDistances[di] = vjds[di]
				}
			}

			if len(commonIDistances) > maxCommon {
				iVector = &vi
				jVector = &vj
				iCommonDistances = commonIDistances
				jCommonDistances = commonJDistances
				maxCommon = len(commonIDistances)
				if maxCommon >= 11 {
					break outer
				}
			}
		}
	}

	if maxCommon < 11 {
		return nil, nil, nil, nil
	}

	return iVector, iCommonDistances, jVector, jCommonDistances
}

func areScannersMerged(scanners []*Scanner) bool {
	count := 0
	for _, s := range scanners {
		if s.merged != nil {
			count++
		}
	}
	return count == 1
}

func main() {

	scanners := getData("../data.txt")

	for i := 0; i < len(scanners); i++ {
		scanners[i].merge(scanners[i].vectorReport, Vector{})
	}

	originScanner := scanners[0]

	for !areScannersMerged(scanners) {
		for i := 1; i < len(scanners); i++ {
			originDistances := originScanner.merged.vectors().distances()
			if scanners[i].merged != nil {
				for t := 0; t < numTransforms; t++ {
					iTransformed := scanners[i].merged.vectors().transform(t)
					iDistances := iTransformed.distances()
					vi, vicds, vj, vjcds := checkOverlap(originDistances, iDistances)
					if vi != nil && vicds != nil && vj != nil && vjcds != nil {
						offset := (*vi).sub(*vj)
						scanners[i].offsetFrom[0] = offsetData{offset: offset, transform: t}
						originScanner.merge(iTransformed, offset)
						scanners[i].merged = nil
						break
					}
				}
			}
		}
	}

	for i := 1; i < len(scanners); i++ {
		fmt.Println(scanners[i].id, scanners[i].offsetFrom[0].offset)
	}

	fmt.Println("num beacons:", len(scanners[0].merged))

	positions := make(Vectors, len(scanners), len(scanners))
	positions[0] = Vector{}
	for i, s := range scanners[1:] {
		positions[i] = s.offsetFrom[0].offset.clone()
	}

	maxMD := 0
	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			md := positions[i].manhattanDistance(positions[j])
			if md > maxMD {
				maxMD = md
			}
		}
	}

	// 11831 too low
	fmt.Println("max manhattan distance between scanners: ", maxMD)

}

func getData(filename string) []*Scanner {
	lines, _ := file.GetLines(filename)
	scanners := make([]*Scanner, 0, 10)

	for i := 0; i < len(lines); i++ {
		s := newScanner(lines[i])

		for i++; i < len(lines) && len(strings.TrimSpace(lines[i])) > 0; i++ {
			v := newVector(lines[i])
			s.add(v)
		}

		scanners = append(scanners, s)
	}

	return scanners
}
