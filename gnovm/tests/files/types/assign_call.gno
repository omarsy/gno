package main

func foo() (int, float32) {
	return 1, 1.0
}

func main() {
	var s, s1 string
	s, s1 = foo()
	println(s)
}

// Error:
// main/assign_call.gno:9:2-15: cannot use int as string

// TypeCheckError:
// main/assign_call.gno:9:10: cannot use foo() (value of type int) as string value in assignment; main/assign_call.gno:9:10: cannot use foo() (value of type float32) as string value in assignment; main/assign_call.gno:8:9: declared and not used: s1
