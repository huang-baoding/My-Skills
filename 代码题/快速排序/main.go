package main

import (
	"fmt"
)

func quickSort(nums []int) {
	if len(nums) < 2 {
		return
	}

	left, right := 0, len(nums)-1 //left是记录有几个比nums[right]小
	tmp := nums[right]
	for i := range nums {
		if nums[i] < tmp {
			nums[i], nums[left] = nums[left], nums[i]
			left++
		}
	}
	//有left个数比nums[right]小，那么把基准数放在nums[left]位置上，左边的数都比基准数小，右边的数都大于等于基准数。
	nums[left], nums[right] = nums[right], nums[left]
	quickSort(nums[:left])
	quickSort(nums[left+1:])
}

func main() {
	nums := []int{1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 0}
	fmt.Println("Original array:", nums)
	quickSort(nums)
	fmt.Println("Sorted array:", nums)
}
