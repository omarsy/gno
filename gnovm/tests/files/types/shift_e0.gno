package main

func main() {
	println(1 << -1)
}

// Error:
// main/shift_e0.gno:4:10-17: invalid operation: negative shift count: (const (-1 <untyped> bigint))

// TypeCheckError:
// main/shift_e0.gno:4:15: invalid operation: negative shift count -1 (untyped int constant)
