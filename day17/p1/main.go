package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reNums = regexp.MustCompile(`(-?\d+)`)
)

func stepX(x int) int {
	if x != 0 {
		if x < 0 {
			x++
		} else {
			x--
		}
	}
	return x
}

func stepY(y int) int {
	return y - 1
}

func getXStepRange(bb *geom.BoundingBox) (int, int) {
	var vxMin, vxMax int
	bbYMax := bb.YMax()
	p := geom.Pos{X: 1, Y: bbYMax}
	d := bb.GetDirection(p)
	v := 1
next:
	for {
		for tv, x := v, v; tv > 0; {
			tv = stepX(tv)
			x += tv
			nd := bb.GetDirection(geom.Pos{X: x, Y: bbYMax})
			if nd != d {
				if nd == d.Opposite() {
					break next
				} else {
					if vxMin == 0 {
						vxMin = v
					}
					vxMax = v
					break
				}
			}
		}
		v++
	}
	return vxMin, vxMax
}

func getYStepRange(bb *geom.BoundingBox) (int, int) {
	var vyMin, vyMax int
	bbXMin := bb.XMin()
	p := geom.Pos{X: bbXMin, Y: 1}
	d := bb.GetDirection(p)
	v := 1
next:
	for {
		for tv, y := v, v; ; {
			tv = stepY(tv)
			y += tv
			nd := bb.GetDirection(geom.Pos{X: bbXMin, Y: y})
			if nd != d {
				if nd == d.Opposite() {
					break next
				} else {
					if vyMin == 0 {
						vyMin = v
					}
					vyMax = v
					break
				}
			}
		}
		v++
	}
	return vyMin, vyMax
}

func main() {
	bb := getData("../data.txt")
	vxMin, vxMax := getXStepRange(bb)
	vyMin, vyMax := getYStepRange(bb)
	fmt.Println(vyMin)

	yMax := 0

	for vx := vxMin; vx <= vxMax; vx++ {
		for yv := -10; yv <= vyMax+50; yv++ {
			p := geom.Pos{X: 0, Y: 0}
			var ty int
			velocity := geom.Pos{X: vx, Y: yv}
			cd := bb.GetDirection(p)
			for !cd.Is([]geom.Direction{geom.South, geom.East}) {
				p.X += velocity.X
				p.Y += velocity.Y
				if p.Y > ty {
					ty = p.Y
				}
				velocity.X = stepX(velocity.X)
				velocity.Y = stepY(velocity.Y)
				cd = bb.GetDirection(p)
				if cd == geom.Unknown {
					if ty > yMax {
						yMax = ty
					}
				}
			}
		}
	}

	fmt.Println(yMax)

}

func getData(filename string) *geom.BoundingBox {
	lines, _ := file.GetLines(filename)

	bb := geom.BoundingBox{}

	matches := reNums.FindAllStringSubmatch(lines[0], -1)

	x1, _ := strconv.ParseInt(matches[0][0], 10, 32)
	y1, _ := strconv.ParseInt(matches[3][0], 10, 32)

	x2, _ := strconv.ParseInt(matches[1][0], 10, 32)
	y2, _ := strconv.ParseInt(matches[2][0], 10, 32)

	bb.SetExtents(int(x1), int(y2), 0, int(x2), int(y1), 0)

	return &bb
}
