package main

var (
	A, B = C, D
	C    = B
	D    = 3
)

func main() {
	println(A, B, C)
}

// Output:
// 3 3 3
