package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

func getCounts(data []string) ([]int, []int) {
	zeros := make([]int, len(data[0]), len(data[0]))
	ones := make([]int, len(data[0]), len(data[0]))
	for _, input := range data {
		for i, c := range input {
			if c == '1' {
				ones[i]++
			} else {
				zeros[i]++
			}
		}
	}
	return ones, zeros
}

func main() {
	data := getData()

	ones, zeros := getCounts(data)

	gamma := make([]byte,len(data[0]),len(data[0]))
	epsilon := make([]byte,len(data[0]),len(data[0]))

	for i := range data[0] {
		if ones[i] > zeros[i] {
			gamma[i], epsilon[i] = '1', '0'
		} else {
			gamma[i], epsilon[i] = '0', '1'
		}
	}

	g, _ := strconv.ParseInt(string(gamma),2,64)
	e, _ := strconv.ParseInt(string(epsilon),2,64)

	fmt.Println(g*e)

	nums := make([]string,0,len(data))
	for i := range data {
		nums = append(nums,data[i])
	}

	for i := range data[0] {
		ones, zeros = getCounts(nums)
		bit := '1'
		if zeros[i] > ones[i] {
			bit = '0'
		}
		tmp := make([]string,0,len(nums))
		for j := range nums {
			if byte(nums[j][i]) == byte(bit) {
				tmp = append(tmp,nums[j])
			}
		}
		nums = tmp
		if len(nums) == 1 {
			break
		}
	}

	ogr, _ := strconv.ParseInt(nums[0],2,64)

	nums = make([]string,0,len(data))
	for i := range data {
		nums = append(nums,data[i])
	}

	for i := range data[0] {
		ones, zeros = getCounts(nums)
		bit := '0'
		if ones[i] < zeros[i] {
			bit = '1'
		}
		tmp := make([]string,0,len(nums))
		for j := range nums {
			if byte(nums[j][i]) == byte(bit) {
				tmp = append(tmp,nums[j])
			}
		}
		nums = tmp
		if len(nums) == 1 {
			break
		}
	}

	co2sr, _ := strconv.ParseInt(nums[0],2,64)

	fmt.Println(ogr * co2sr)

}

func getData() []string {
	lines, _ := file.GetLines("../data.txt")
	return lines
}
