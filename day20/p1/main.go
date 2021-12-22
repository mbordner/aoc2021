package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"fmt"
	"strconv"
	"strings"
)

const (
	pixelON  = '#'
	pixelOFF = '.'
)

type ImageEnhancementAlgorithm []byte

type Image struct {
	bb     geom.BoundingBox
	bg     byte
	pixels map[geom.Pos]byte
}

func newImage() *Image {
	i := new(Image)
	i.bg = pixelOFF
	i.pixels = make(map[geom.Pos]byte)

	return i
}

func (i *Image) print() {
	xOffset := i.bb.XMin() * -1
	yOffset := i.bb.YMin() * -1
	positions := i.bb.GetPositions()
	positions = positions.Transform(xOffset, yOffset, 0)
	var x, y int
	lineWidth := i.bb.XMax() - i.bb.XMin() + 1
	line := make([]byte, lineWidth, lineWidth)
	for p := 0; p < len(positions); p++ {

		line[x] = i.getPixel(positions[p].Transform(-xOffset, -yOffset, 0))
		x++

		if x == lineWidth {
			fmt.Println(string(line))
			line = make([]byte, lineWidth, lineWidth)
			y++
			x = 0
		}

	}
}

func (i *Image) setPixel(pos geom.Pos, val byte) {
	i.bb.Extend(pos)
	i.pixels[pos] = val
}

func (i *Image) expandSpace(amt int, pixel byte) {
	i.bg = pixel
	xmin, xmax, ymin, ymax := i.bb.XMin(), i.bb.XMax(), i.bb.YMin(), i.bb.YMax()

	// top 2 rows beyond y min
	for y := ymin - amt; y < ymin; y++ {
		for x := xmin - amt; x < xmax+amt+1; x++ {
			i.setPixel(geom.Pos{Y: y, X: x}, i.bg)
		}
	}

	// bottom 2 rows beyond y max
	for y := ymax + 1; y < ymax+amt+1; y++ {
		for x := xmin - amt; x < xmax+amt+1; x++ {
			i.setPixel(geom.Pos{Y: y, X: x}, i.bg)
		}
	}

	// left two columns
	for y := ymin; y < ymax+1; y++ {
		for x := xmin - amt; x < xmin; x++ {
			i.setPixel(geom.Pos{Y: y, X: x}, i.bg)
		}
	}

	// right 2 columns
	for y := ymin; y < ymax+1; y++ {
		for x := xmax + 1; x < xmax+amt+1; x++ {
			i.setPixel(geom.Pos{Y: y, X: x}, i.bg)
		}
	}
}

func (i *Image) getBGPixel() byte {
	return i.bg
}

func (i *Image) getPixel(pos geom.Pos) byte {
	if i.bb.Contains(pos) {
		if v, exists := i.pixels[pos]; exists {
			return v
		}
		return pixelOFF
	}
	return i.getBGPixel()
}

func (i *Image) getSurroundingValues(pos geom.Pos) []byte {
	vals := make([]byte, 0, 9)
	for y := pos.Y - 1; y <= pos.Y+1; y++ {
		for x := pos.X - 1; x <= pos.X+1; x++ {
			vals = append(vals, i.getPixel(geom.Pos{Y: y, X: x}))
		}
	}
	return vals
}

func (i *Image) getPixels() []geom.Pos {
	return i.bb.GetPositions()
}

func (i *Image) getOnPixelCount() int {
	count := 0
	for _, v := range i.pixels {
		if v == pixelON {
			count++
		}
	}
	return count
}

func getAlgorithmIndex(vals []byte) int {
	digits := make([]byte, len(vals), len(vals))
	for i := range vals {
		if vals[i] == pixelON {
			digits[i] = '1'
		} else {
			digits[i] = '0'
		}
	}
	val, _ := strconv.ParseInt(string(digits), 2, 32)
	return int(val)
}

func main() {
	testVal := getAlgorithmIndex([]byte(`...#...#.`))
	if testVal != 34 {
		panic("get algorithm index doesn't return expected values")
	}

	images := make([]*Image, 0, 100)

	imgEnhanceAlgo, initImg := getData("../data.txt")

	expandEven := imgEnhanceAlgo[getAlgorithmIndex([]byte(strings.Repeat(".", 9)))]
	expandOdd := imgEnhanceAlgo[getAlgorithmIndex([]byte(strings.Repeat(string(expandEven), 9)))]
	backgrounds := []byte{expandEven, expandOdd}

	images = append(images, initImg)

	numApplications := 0

	numExpand := 1

	initImg.expandSpace(numExpand, pixelOFF)

	for numApplications < 50 {

		curImg := images[len(images)-1]
		newImg := newImage()

		for _, p := range curImg.getPixels() {
			sv := curImg.getSurroundingValues(p)
			index := getAlgorithmIndex(sv)
			newImg.setPixel(p, imgEnhanceAlgo[index])
		}

		newImg.expandSpace(numExpand, backgrounds[numApplications%2])

		images = append(images, newImg)

		numApplications++
	}

	fmt.Println("\n---\n")
	for _, img := range images {
		img.print()
		fmt.Println("\n---\n")
	}

	fmt.Println("last image on pixel count: ", images[len(images)-1].getOnPixelCount())

}

func getData(filename string) (ImageEnhancementAlgorithm, *Image) {
	lines, _ := file.GetLines(filename)
	iea := ImageEnhancementAlgorithm(lines[0])
	image := newImage()

	for y, line := range lines[2:] {
		for x := 0; x < len(line); x++ {
			image.setPixel(geom.Pos{X: x, Y: y}, line[x])
		}
	}

	return iea, image
}
