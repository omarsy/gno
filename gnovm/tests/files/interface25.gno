package main

import "fmt"

func main() {
	m := make(map[string]any)
	m["A"] = 1
	for _, v := range m {
		fmt.Println(v)
	}
}

// Output:
// 1
