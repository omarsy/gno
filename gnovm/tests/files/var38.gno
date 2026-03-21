package main

// True circular dependency with multi-value init should still be caught.
var (
	A, B = C, D
	C    = A
	D    = 3
)

func main() {
	println(A, B, C)
}

// Error:
// main/var38.gno:1:1-12:2: invalid recursive value: A -> C -> A

// TypeCheckError:
// main/var38.gno:5:2: initialization cycle for A; main/var38.gno:5:2: 	A refers to C; main/var38.gno:6:2: 	C refers to A
