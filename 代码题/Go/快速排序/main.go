package main

import "fmt"

var nums []int = []int{2, 34, 55, 2, 45, 65, 76, 9, 28}

func main() {
	quickSort(nums)
	fmt.Println(nums)
}

func quickSort(nums []int) {
	if len(nums) <= 1 {
		return
	}

	tmp := nums[0]
	left := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] < tmp {
			nums[left], nums[i] = nums[i], nums[left]
			left++
		}
	}
	quickSort(nums[:left])
	quickSort(nums[left+1:])
}
