package main

import (
	"strconv"
)

type Error1 int64

func (e Error1) Error() string {
	return "error: " + strconv.Itoa(int(e))
}

type Error2 int

func (e Error2) Error() string {
	return "error: " + strconv.Itoa(int(e))
}

// both not const, typed
func main() {
	var e1 Error1 = Error1(0)
	var e2 Error2 = Error2(0)
	e1 += e2
	println(e1)
}

// Error:
// main/add_assign_e0.gno:23:2-10: invalid operation: mismatched types main.Error1 and main.Error2

// TypeCheckError:
// main/add_assign_e0.gno:23:2: invalid operation: e1 += e2 (mismatched types Error1 and Error2)
