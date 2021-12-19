package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParsingSFNum(t *testing.T) {

	tests := []struct {
		Input string
		Depth int
	}{
		{Input: `[[[[[9,8],1],2],3],4]`, Depth: 5},
		{Input: `[[[[1,2],[3,4]],[[5,6],[7,8]]],9]`, Depth: 4},
		{Input: `[[1,2],3]`, Depth: 2},
		{Input: `[1,2]`, Depth: 1},
		{Input: `[[1,9],[8,5]]`, Depth: 2},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input)

			assert.NotNil(t, sfn)
			assert.Equal(t, tc.Depth, sfn.getDepth())
			assert.Equal(t, tc.Input, sfn.String())

		})
	}
}

func Test_Magnitude(t *testing.T) {

	tests := []struct {
		Input     string
		Magnitude int
	}{
		{Input: `[[1,2],[[3,4],5]]`, Magnitude: 143},
		{Input: `[[[[0,7],4],[[7,8],[6,0]]],[8,1]]`, Magnitude: 1384},
		{Input: `[[[[1,1],[2,2]],[3,3]],[4,4]]`, Magnitude: 445},
		{Input: `[[[[3,0],[5,3]],[4,4]],[5,5]]`, Magnitude: 791},
		{Input: `[[[[5,0],[7,4]],[5,5]],[6,6]]`, Magnitude: 1137},
		{Input: `[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]`, Magnitude: 3488},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input)

			assert.NotNil(t, sfn)
			assert.Equal(t, tc.Magnitude, sfn.getMagnitude())
			assert.Equal(t, tc.Input, sfn.String())

		})
	}
}

func Test_GetLeafs(t *testing.T) {
	tests := []struct {
		Input  string
		Values []int
	}{
		{Input: `[[[[[9,8],1],2],3],4]`, Values: []int{9, 8, 1, 2, 3, 4}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input)

			assert.NotNil(t, sfn)

			leafs := sfn.getLeafs()
			assert.Equal(t, len(tc.Values), len(leafs))

			for i, tsfn := range leafs {
				assert.True(t, tsfn.isLeaf())
				assert.Equal(t, tc.Values[i], tsfn.getValue())
			}

		})
	}
}

func Test_ToExplode(t *testing.T) {

	tests := []struct {
		Input string
		Size  int
	}{
		{Input: `[[[[[9,8],1],2],3],4]`, Size: 1},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input)

			sfns := sfn.toExplode(1)

			assert.NotNil(t, sfns)
			assert.Equal(t, tc.Size, len(sfns))
			assert.True(t, sfns[0].left.isLeaf())
			assert.True(t, sfns[0].right.isLeaf())

		})
	}

}

func Test_ToSplit(t *testing.T) {
	tests := []struct {
		Input string
		Size  int
	}{
		{Input: `[[[[[10,8],1],2],3],4]`, Size: 1},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input)

			sfns := sfn.toSplit()

			assert.NotNil(t, sfns)
			assert.Equal(t, tc.Size, len(sfns))
			assert.True(t, sfns[0].isLeaf())

		})
	}
}

func Test_ToAdd(t *testing.T) {
	tests := []struct {
		Input1        string
		Input2        string
		AfterAddition string
	}{
		{Input1: `[[[[4,3],4],4],[7,[[8,4],9]]]`, Input2: `[1,1]`, AfterAddition: `[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]`},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input1)
			assert.NotNil(t, sfn)

			osfn := newSFNum(tc.Input2)
			assert.NotNil(t, osfn)

			sfn = sfn.add(osfn)
			assert.Equal(t, tc.AfterAddition, sfn.String())

		})
	}
}

func Test_AddReduce(t *testing.T) {
	tests := []struct {
		Input1      string
		Input2      string
		AfterReduce string
	}{
		{Input1: `[[[[4,3],4],4],[7,[[8,4],9]]]`, Input2: `[1,1]`, AfterReduce: `[[[[0,7],4],[[7,8],[6,0]]],[8,1]]`},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			sfn := newSFNum(tc.Input1)
			assert.NotNil(t, sfn)

			osfn := newSFNum(tc.Input2)
			assert.NotNil(t, osfn)

			sfn = sfn.addAndReduce(osfn)
			reduced := sfn.String()
			assert.Equal(t, tc.AfterReduce, reduced)

		})
	}
}

func Test_AddList(t *testing.T) {
	list := []string{
		`[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]`,
		`[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]`,
		`[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]`,
		`[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]`,
		`[7,[5,[[3,8],[1,4]]]]`,
		`[[2,[2,2]],[8,[8,1]]]`,
		`[2,9]`,
		`[1,[[[9,3],9],[[9,0],[0,7]]]]`,
		`[[[5,[7,4]],7],1]`,
		`[[[[4,2],2],6],[8,7]]`,
	}

	sfn := newSFNum(list[0])
	for _, num := range list[1:] {
		sfn = sfn.addAndReduce(newSFNum(num))
	}

	assert.Equal(t, `[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]`, sfn.String())
}

func Test_Sample(t *testing.T) {
	list := getList([]string{
		`[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]`,
		`[[[5,[2,8]],4],[5,[[9,9],0]]]`,
		`[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]`,
		`[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]`,
		`[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]`,
		`[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]`,
		`[[[[5,4],[7,7]],8],[[8,3],8]]`,
		`[[9,3],[[9,9],[6,[4,9]]]]`,
		`[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]`,
		`[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]`,
	})

	sfn := list[0]
	for _, num := range list[1:] {
		sfn = sfn.addAndReduce(num)
	}

	assert.Equal(t, `[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]`, sfn.String())
	assert.Equal(t, 4140, sfn.getMagnitude())
}

func Test_LargestMagnitude(t *testing.T) {
	lines := []string{
		`[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]`,
		`[[[5,[2,8]],4],[5,[[9,9],0]]]`,
		`[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]`,
		`[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]`,
		`[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]`,
		`[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]`,
		`[[[[5,4],[7,7]],8],[[8,3],8]]`,
		`[[9,3],[[9,9],[6,[4,9]]]]`,
		`[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]`,
		`[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]`,
	}

	max := getLargestMagnitudeAddingOnlyTwo(lines)

	assert.Equal(t, 3993, max)
}
