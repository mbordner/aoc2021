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
	MinX int
	MaxX int
	MinY int
	MaxY int
	MinZ int
	MaxZ int
}

func (bb *BoundingBox) SetExtents(x1, y1, z1, x2, y2, z2 int) {
	bb.MinX = x1
	bb.MinY = y1
	bb.MinZ = z1
	bb.MaxX = x2
	bb.MaxY = y2
	bb.MaxZ = z2
}

func (bb BoundingBox) XMin() int {
	return bb.MinX
}

func (bb BoundingBox) XMax() int {
	return bb.MaxX
}

func (bb BoundingBox) YMin() int {
	return bb.MinY
}

func (bb BoundingBox) YMax() int {
	return bb.MaxY
}

func (bb BoundingBox) ZMin() int {
	return bb.MinZ
}

func (bb BoundingBox) ZMax() int {
	return bb.MaxZ
}

func (bb BoundingBox) String() string {
	p1 := Pos{X: bb.MinX, Y: bb.MinY}
	p2 := Pos{X: bb.MaxX, Y: bb.MaxY}
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

func (bb *BoundingBox) Extend(p Pos) {
	if p.X < bb.MinX {
		bb.MinX = p.X
	}
	if p.X > bb.MaxX {
		bb.MaxX = p.X
	}
	if p.Y > bb.MaxY {
		bb.MaxY = p.Y
	}
	if p.Y < bb.MinY {
		bb.MinY = p.Y
	}
	if p.Z < bb.MinZ {
		bb.MinZ = p.Z
	}
	if p.Z > bb.MaxZ {
		bb.MaxZ = p.Z
	}
}

func (bb *BoundingBox) Contains(p Pos) bool {
	if p.X < bb.MinX || p.X > bb.MaxX {
		return false
	}
	if p.Y < bb.MinY || p.Y > bb.MaxY {
		return false
	}
	if p.Z < bb.MinZ || p.Z > bb.MaxZ {
		return false
	}
	return true
}

func (bb *BoundingBox) Surrounds(obb *BoundingBox) bool {
	if obb.MinX < bb.MinX {
		return false
	}
	if obb.MaxX > bb.MaxX {
		return false
	}
	if obb.MinY < bb.MinY {
		return false
	}
	if obb.MaxY > bb.MaxY {
		return false
	}
	if obb.MinZ < bb.MinZ {
		return false
	}
	if obb.MaxZ > bb.MaxZ {
		return false
	}
	return true
}

func (bb *BoundingBox) GetDirection(p Pos) Direction {
	dir := 0
	if p.X < bb.MinX {
		dir |= int(West)
	}
	if p.X > bb.MaxX {
		dir |= int(East)
	}
	if p.Y > bb.MaxY {
		dir |= int(North)
	}
	if p.Y < bb.MinY {
		dir |= int(South)
	}
	return Direction(dir)
}

func (bb *BoundingBox) Intersects(p1, p2 Pos) bool {
	return false
}

func (bb *BoundingBox) DistanceFromEdge(p Pos) int {
	d := math.MaxInt64

	t := bb.MaxX - p.X
	if t < d {
		d = t
	}

	t = p.X - bb.MinX
	if t < d {
		d = t
	}

	t = p.Y - bb.MinY
	if t < d {
		d = t
	}

	t = bb.MaxY - p.Y
	if t < d {
		d = t
	}

	return d
}

type Positions []Pos

type Pos struct {
	X int
	Y int
	Z int
}

func (p Pos) Transform(x, y, z int) Pos {
	return Pos{X: p.X + x, Y: p.Y + y, Z: p.Z + z}
}

func (p Pos) Clone() Pos {
	return Pos{X: p.X, Y: p.Y, Z: p.Z}
}

func (p Pos) String() string {
	return fmt.Sprintf("{x:%d, y:%d, z:%d}", p.X, p.Y, p.Z)
}

func Abs(n int) int {
	if n > 0 {
		return n
	}
	return -n
}

func (bb *BoundingBox) GetPositionsSize() uint64 {
	xs := Abs(bb.MaxX-bb.MinX) + 1
	ys := Abs(bb.MaxY-bb.MinY) + 1
	zs := Abs(bb.MaxZ-bb.MinZ) + 1
	return uint64(xs) * uint64(ys) * uint64(zs)
}

func (bb *BoundingBox) GetPositions() Positions {
	poss := make(Positions, 0, ((bb.MaxX-bb.MinX)+1)*((bb.MaxY-bb.MinY)+1*((bb.MaxZ-bb.MinZ)+1)))
	for z := bb.MinZ; z <= bb.MaxZ; z++ {
		for y := bb.MinY; y <= bb.MaxY; y++ {
			for x := bb.MinX; x <= bb.MaxX; x++ {
				poss = append(poss, Pos{Z: z, Y: y, X: x})
			}
		}
	}
	return poss
}

func (ps *Positions) Transform(x, y, z int) Positions {
	for i := 0; i < len(*ps); i++ {
		(*ps)[i].X += x
		(*ps)[i].Y += y
		(*ps)[i].Z += z
	}
	return *ps
}
