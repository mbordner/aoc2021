package array

import (
	"strconv"
	"strings"
)

func ToIntArray(s string) []int {
	tokens := strings.Split(s, ",")
	nums := make([]int, len(tokens), len(tokens))
	for i := range tokens {
		val, _ := strconv.ParseInt(tokens[i], 10, 32)
		nums[i] = int(val)
	}
	return nums
}
