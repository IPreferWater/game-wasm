package main

import (
	"fmt"
	"testing"
)

func TestArr(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4}
	aaa(arr)
	t.Error(arr)
}

func aaa(arr []int){
	arr = pop(arr,2)
	fmt.Printf("%v from nested \n",arr)
}

func pop(s []int, i int) []int {
	ret := make([]int, 0)
    ret = append(ret, s[:i]...)
    return append(ret, s[i+1:]...)
}
