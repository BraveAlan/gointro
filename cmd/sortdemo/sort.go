package main

import (
	"fmt"
	"sort"
)

func main() {
	// Creats a slice of int
	a := []int{3, 6, 2, 1, 9}
	sort.Ints(a)
	for _, v := range a {
		fmt.Println(v)
	}
}
