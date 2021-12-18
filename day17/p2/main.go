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

func main() {
	bb := getData("../data.txt")

	yMax := 0
	velocities := make([]geom.Pos, 0, 100)

	for y := bb.YMin(); y < -bb.YMin(); y++ {

		var d geom.Direction
		var curPos geom.Pos
	inner:
		for x := 1; x <= bb.XMax()+1; x++ {
			initialVelocity := geom.Pos{Y: y, X: x}
			//fmt.Println("trying velocity", initialVelocity)
			velocity := geom.Pos{Y: y, X: x}
			curPos = geom.Pos{Y: 0, X: 0}
			curYMax := 0
			for {
				curPos.X += velocity.X
				curPos.Y += velocity.Y

				if curPos.Y > curYMax {
					curYMax = curPos.Y
				}

				if bb.Contains(curPos) {
					velocities = append(velocities, initialVelocity)
					if curYMax > yMax {
						yMax = curYMax
					}
					continue inner
				} else {
					d = bb.GetDirection(curPos)
					if d.Is([]geom.Direction{geom.South, geom.East}) {
						continue inner
					}
				}

				velocity.X = stepX(velocity.X)
				velocity.Y = stepY(velocity.Y)
			}
		}
	}

	fmt.Println(yMax)
	fmt.Println(len(velocities))
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
