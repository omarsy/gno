package main

import (
	"strconv"
)

type Error int8

func (e Error) Error() string {
	return "error: " + strconv.Itoa(int(e))
}

// left is typed const, right is untyped const
// NOTE: overflow
func main() {
	if Error(1) == 128 { // note, this would overflow as expected
		println("what the firetruck?")
	} else {
		println("something else")
	}
}

// Error:
// main/eql_0c2.gno:16:5-20: bigint overflows target kind

// TypeCheckError:
// main/eql_0c2.gno:16:17: 128 (untyped int constant) overflows int8
