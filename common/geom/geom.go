package geom

import (
	"fmt"
	"math"
)

type Direction int

const (
	Unknown Direction = 0
	North   Direction = 1
	South   Direction = 2
	West    Direction = 4
	East    Direction = 8
)

func (d Direction) Is(dirs []Direction) bool {
	for _, dir := range dirs {
		if int(d)&int(dir) == int(dir) {
			return true
		}
	}
	return false
}

func (d Direction) Opposite() Direction {
	o := 0
	if int(d)&int(North) == int(North) {
		o |= int(South)
	}
	if int(d)&int(West) == int(West) {
		o |= int(East)
	}
	if int(d)&int(South) == int(South) {
		o |= int(North)
	}
	if int(d)&int(East) == int(East) {
		o |= int(West)
	}
	return Direction(o)
}

type BoundingBox struct {
	xMin int
	xMax int
	yMin int
	yMax int
	zMin int
	zMax int
}

func (bb *BoundingBox) SetExtents(x1, y1, z1, x2, y2, z2 int) {
	bb.xMin = x1
	bb.yMin = y1
	bb.zMin = z1
	bb.xMax = x2
	bb.yMax = y2
	bb.zMax = z2
}

func (bb BoundingBox) XMin() int {
	return bb.xMin
}

func (bb BoundingBox) XMax() int {
	return bb.xMax
}

func (bb BoundingBox) YMin() int {
	return bb.yMin
}

func (bb BoundingBox) YMax() int {
	return bb.yMax
}

func (bb BoundingBox) ZMin() int {
	return bb.zMin
}

func (bb BoundingBox) ZMax() int {
	return bb.zMax
}

func (bb BoundingBox) String() string {
	p1 := Pos{X: bb.xMin, Y: bb.yMin}
	p2 := Pos{X: bb.xMax, Y: bb.yMax}
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

func (bb *BoundingBox) Extend(p Pos) {
	if p.X < bb.xMin {
		bb.xMin = p.X
	}
	if p.X > bb.xMax {
		bb.xMax = p.X
	}
	if p.Y > bb.yMax {
		bb.yMax = p.Y
	}
	if p.Y < bb.yMin {
		bb.yMin = p.Y
	}
	if p.Z < bb.zMin {
		bb.zMin = p.Z
	}
	if p.Z > bb.zMax {
		bb.zMax = p.Z
	}
}

func (bb *BoundingBox) Contains(p Pos) bool {
	if p.X < bb.xMin || p.X > bb.xMax {
		return false
	}
	if p.Y < bb.yMin || p.Y > bb.yMax {
		return false
	}
	if p.Z < bb.zMin || p.Z > bb.zMax {
		return false
	}
	return true
}

func (bb *BoundingBox) GetDirection(p Pos) Direction {
	dir := 0
	if p.X < bb.xMin {
		dir |= int(West)
	}
	if p.X > bb.xMax {
		dir |= int(East)
	}
	if p.Y > bb.yMax {
		dir |= int(North)
	}
	if p.Y < bb.yMin {
		dir |= int(South)
	}
	return Direction(dir)
}

func (bb *BoundingBox) Intersects(p1, p2 Pos) bool {
	return false
}

func (bb *BoundingBox) DistanceFromEdge(p Pos) int {
	d := math.MaxInt64

	t := bb.xMax - p.X
	if t < d {
		d = t
	}

	t = p.X - bb.xMin
	if t < d {
		d = t
	}

	t = p.Y - bb.yMin
	if t < d {
		d = t
	}

	t = bb.yMax - p.Y
	if t < d {
		d = t
	}

	return d
}

type Pos struct {
	X int
	Y int
	Z int
}

func (p Pos) Clone() Pos {
	return Pos{X: p.X, Y: p.Y, Z: p.Z}
}

func (p Pos) String() string {
	return fmt.Sprintf("{x:%d, y:%d, z:%d}", p.X, p.Y, p.Z)
}

func (bb *BoundingBox) GetPositions() []Pos {
	poss := make([]Pos, 0, ((bb.xMax-bb.xMin)+1)*((bb.yMax-bb.yMin)+1))
	for y := bb.yMin; y <= bb.yMax; y++ {
		for x := bb.xMin; x <= bb.xMax; x++ {
			poss = append(poss, Pos{Y: y, X: x})
		}
	}
	return poss
}
