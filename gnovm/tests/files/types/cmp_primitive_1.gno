package main

type Error string

func (e Error) Error() string {
	return "error: " + string(e)
}

// left is untyped const, right is typed const
// left is not assignable to right
// a) it's (untyped) bigint
// b) base type of right is string
func main() {
	if 1 == Error(1) {
		println("what the firetruck?")
	} else {
		println("something else")
	}
}

// Error:
// main/cmp_primitive_1.gno:14:5-18: cannot use untyped Bigint as StringKind

// TypeCheckError:
// main/cmp_primitive_1.gno:14:10: invalid operation: 1 == Error(1) (mismatched types untyped int and Error)
