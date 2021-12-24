package geom

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Volume(t *testing.T) {

	cases := []struct {
		input  string
		volume uint64
	}{
		{input: `0,0,0,1,1,1`, volume: uint64(1)},
		{input: `0,0,0,0,0,0`, volume: uint64(0)},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			v := NewCuboid(tc.input).Volume()
			assert.Equal(t, tc.volume, v)
		})
	}
}

func Test_CubuoidsVolume(t *testing.T) {

	cases := []struct {
		cube     string
		splitpt  string
		expected uint64
	}{
		{cube: `0,0,0,5,5,5`, splitpt: `2,2,2`, expected: 64},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.cube)
			splits := c.SplitAt(NewPoint(tc.splitpt))
			assert.Equal(t, tc.expected, splits.Volume())
		})
	}

}

func Test_PointsCount(t *testing.T) {

	cases := []struct {
		input  string
		volume uint64
	}{
		{input: `0,0,0,1,1,1`, volume: uint64(8)},
		{input: `0,0,0,0,0,0`, volume: uint64(1)},
		{input: `0,0,0,2,2,2`, volume: uint64(27)},
		{input: `0,0,0,2,2,1`, volume: uint64(18)},
		{input: `0,0,0,2,2,0`, volume: uint64(9)},
		{input: `0,0,0,3,3,0`, volume: uint64(16)},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			v := NewCuboid(tc.input).PointsCount()
			assert.Equal(t, tc.volume, v)
		})
	}
}

func Test_PointsContainsViaCorners(t *testing.T) {
	cases := []struct {
		cuboid   string
		points   []string
		expected bool
	}{
		{cuboid: `0,0,0,1,1,1`, points: []string{`0,0,0`, `0,1,0`, `1,1,0`, `1,0,0`, `0,0,1`, `0,1,1`, `1,1,1`, `1,0,1`}, expected: true},
		{cuboid: `0,0,0,1,1,1`, points: []string{`2,2,2`}, expected: false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.cuboid)
			corners := c.Corners()
			for _, s := range tc.points {
				p := NewPoint(s)
				val := corners.Contains(p)

				assert.Equal(t, tc.expected, val)
			}

		})
	}

}

func Test_TransformPoint(t *testing.T) {
	cases := []struct {
		point    string
		vector   string
		expected string
	}{
		{point: `0,0,0`, vector: `1,-1,1`, expected: `1,-1,1`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewPoint(tc.point).Transform(NewVector(tc.vector)).String())
		})
	}
}

func Test_CuboidContains(t *testing.T) {
	cases := []struct {
		cuboid   string
		point    string
		expected bool
	}{
		{`0,0,0,2,2,2`, `1,1,1`, true},
		{`0,0,0,2,2,2`, `-1,-1,-1`, false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewCuboid(tc.cuboid).Contains(NewPoint(tc.point)))
		})
	}
}

func Test_CuboidsContains(t *testing.T) {
	cases := []struct {
		cuboids  []string
		check    string
		expected bool
	}{
		{
			cuboids:  []string{`0,0,0,1,1,1`, `2,2,2,3,3,3`, `4,4,4,5,5,5`},
			check:    `2,2,2,3,3,3`,
			expected: true,
		},
		{
			cuboids:  []string{`0,0,0,1,1,1`, `2,2,2,3,3,3`, `4,4,4,5,5,5`},
			check:    `6,6,6,7,7,7`,
			expected: false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			cuboids := make(Cuboids, 0, len(tc.cuboids))
			for _, c := range tc.cuboids {
				cuboids = append(cuboids, NewCuboid(c))
			}
			assert.Equal(t, tc.expected, cuboids.Contains(NewCuboid(tc.check)))
		})
	}
}

func Test_IntersectingCorners(t *testing.T) {
	cases := []struct {
		c1      string
		c2      string
		corners []string
	}{
		{`0,0,0,2,2,2`, `-1,-1,-1,1,1,1`, []string{`1,1,1`}},
		{`0,0,0,2,2,2`, `-2,-2,-2,-1,-1,-1`, []string{}},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			corners := NewCuboid(tc.c1).IntersectingCorners(NewCuboid(tc.c2))

			assert.Equal(t, len(tc.corners), len(corners))
		})
	}
}

func Test_IsCorner(t *testing.T) {
	cases := []struct {
		c1       string
		p        string
		expected bool
	}{
		{`0,0,0,1,1,1`, `1,1,1`, true},
		{`0,0,0,1,1,1`, `-1,-1,-1`, false},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewCuboid(tc.c1).IsCorner(NewPoint(tc.p)))
		})
	}
}

func Test_CuboidTransform(t *testing.T) {
	cases := []struct {
		c1       string
		v        string
		expected string
	}{
		{`0,0,0,1,1,1`, `1,1,1`, `1,1,1,2,2,2`},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			assert.Equal(t, tc.expected, NewCuboid(tc.c1).Transform(NewVector(tc.v)).String())
		})
	}
}

func Test_SplitAt(t *testing.T) {
	type pointCheck struct {
		p     string
		count int
	}

	cases := []struct {
		c      string
		p      string
		l      int
		checks []pointCheck // loops on points across the split cuboids and assets the points exit in *count* cuboids
	}{
		{
			c: `-2,-2,-2,2,2,2`,
			p: `0,0,0`,
			l: 8,
		},
		{
			c: `0,0,0,2,2,2`,
			p: `1,1,1`,
			l: 8,
		},
		{
			c: `0,0,0,3,3,3`,
			p: `1,1,1`,
			l: 8,
		},
		{
			c: `0,0,0,4,4,4`,
			p: `1,1,1`,
			l: 8,
		},
		{
			c: `0,0,0,4,4,4`,
			p: `2,2,2`,
			l: 8,
		},
		{
			c: `0,0,0,5,5,5`,
			p: `2,2,2`,
			l: 8,
		},
		{
			c: `0,0,0,10,10,10`,
			p: `2,2,2`,
			l: 8,
			checks: []pointCheck{
				{
					p:     `2,2,2`,
					count: 1,
				},
			},
		},
		{
			c: `0,0,0,10,10,10`,
			p: `-1,-1,-1`,
			l: 0,
		},
		{
			c: `0,0,0,10,10,10`,
			p: `0,0,0`,
			l: 8,
			checks: []pointCheck{
				{
					p:     `0,0,0`,
					count: 1,
				},
				{
					p:     `1,1,1`,
					count: 1,
				},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c := NewCuboid(tc.c)
			p := NewPoint(tc.p)

			splits := c.SplitAt(p)

			assert.Equal(t, tc.l, len(splits))

			if len(splits) > 0 {
				cuboidPointsCount := c.PointsCount()
				splitsPointsCount := splits.PointsCount()
				assert.Equal(t, cuboidPointsCount, splitsPointsCount)
			}

			if tc.checks != nil && len(tc.checks) > 0 {
				for _, check := range tc.checks {
					tp := NewPoint(check.p)

					count := 0
					for _, c := range splits {
						if c.Contains(tp) {
							count++
						}
					}

					assert.Equal(t, check.count, count)
				}
			}
		})
	}
}

func Test_IntersectCuboids(t *testing.T) {

	type results struct {
		c1   []string
		both []string
		c2   []string
	}

	cases := []struct {
		c1     string
		c2     string
		checks results
	}{
		{
			c1: `5,5,5,10,10,10`,
			c2: `0,0,0,4,4,4`,
			checks: results{
				c1:   []string{`5,5,5,10,10,10`},
				both: []string{},
				c2:   []string{`0,0,0,4,4,4`},
			},
		},
		{
			c1: `0,0,0,5,5,5`,
			c2: `-5,-5,-5,2,2,2`,
			checks: results{
				c1:   []string{`3,0,0,5,5,5`, `0,3,0,2,5,5`},
				both: []string{`0,0,0,2,2,2`},
				c2:   []string{`-5,-5,-5,-1,2,2`, `0,-5,-5,2,-1,2`},
			},
		},
		{
			c1: `0,0,0,10,10,10`,
			c2: `0,0,0,10,10,5`,
			checks: results{
				c1:   []string{`0,0,6,10,10,10`},
				both: []string{`0,0,0,10,10,5`},
				c2:   []string{},
			},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			c1 := NewCuboid(tc.c1)
			c2 := NewCuboid(tc.c2)
			fromC1, both, fromC2 := c1.Intersect(c2)

			c1pc := c1.PointsCount()
			c2pc := c2.PointsCount()
			interpc := both.PointsCount()

			fromc1pc := fromC1.PointsCount()
			fromc2pc := fromC2.PointsCount()

			val1 := c1pc + c2pc - interpc
			val2 := fromc1pc + interpc + fromc2pc

			assert.Equal(t, val1, val2)

			checkFromC1 := make(Cuboids, 0, len(tc.checks.c1))
			for _, c := range tc.checks.c1 {
				checkFromC1 = append(checkFromC1, NewCuboid(c))
			}

			checkBoth := make(Cuboids, 0, len(tc.checks.both))
			for _, c := range tc.checks.both {
				checkBoth = append(checkBoth, NewCuboid(c))
			}

			checkFromC2 := make(Cuboids, 0, len(tc.checks.c2))
			for _, c := range tc.checks.c2 {
				checkFromC2 = append(checkFromC2, NewCuboid(c))
			}

			assert.Equal(t, len(checkFromC1), len(fromC1))
			assert.Equal(t, len(checkBoth), len(both))
			assert.Equal(t, len(checkFromC2), len(fromC2))

			for _, check := range [][]Cuboids{{checkFromC1, fromC1}, {checkBoth, both}, {checkFromC2, fromC2}} {
				for _, c := range check[0] {
					assert.True(t, check[1].Contains(c))
				}
			}
		})
	}
}
