package main

import "fmt"

func unpack_nums(nums ...int) {
	fmt.Println(nums)
}

func Unpack_Slice_Into_Individual_Args() {
	arr := []int{10, 20, 30}
	unpack_nums(arr...)
}
