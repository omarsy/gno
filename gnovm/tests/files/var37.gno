package main

// Multi-value init where second var depends on first var's dependency.
var (
	A, B = 1, C
	C    = A
)

func main() {
	println(A, B, C)
}

// Output:
// 1 1 1
