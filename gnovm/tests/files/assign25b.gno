package main

func foo() (int, bool) {
	return 1, true
}

func main() {
	a, b, c := 2, 3, 4, foo()

	println(a, b, c)
}

// Error:
// main/assign25b.gno:8:2-27: assignment mismatch: 3 variable(s) but 4 value(s)

// TypeCheckError:
// main/assign25b.gno:8:13: assignment mismatch: 3 variables but 4 values
