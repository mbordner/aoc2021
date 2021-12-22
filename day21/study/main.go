package main

import "fmt"

func main() {
	vals := []int{1, 2, 3}

	perms := make([][]int, 0, 27)

	sums := make(map[int]int)
	for i := 0; i < len(vals); i++ {
		for j := 0; j < len(vals); j++ {
			for k := 0; k < len(vals); k++ {
				sum := vals[i] + vals[j] + vals[k]
				//str := fmt.Sprintf("{%d,%d,%d}",vals[i],vals[j],vals[k])

				if c, e := sums[sum]; e {
					sums[sum] = c + 1
				} else {
					sums[sum] = 1
				}

				perms = append(perms, []int{vals[i], vals[j], vals[k]})
			}
		}
	}

	fmt.Println(perms)
	fmt.Println(sums)

}

// Perm calls f with each permutation of a.
func Perm(a []rune, f func([]rune)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []rune, f func([]rune), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}
