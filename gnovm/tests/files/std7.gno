package main

import (
	"std"

	"gno.land/p/demo/testutils"
)

func inner() {
	caller1 := std.CallerAt(1)
	println(caller1)
	caller2 := std.CallerAt(2)
	println(caller2)
	caller3 := std.CallerAt(3)
	println(caller3)
}

func main() {
	testutils.WrapCall(inner)
}

// Output:
// g17rgsdnfxzza0sdfsdma37sdwxagsz378833ca4
// g1ssrgtfce6pzs8tp7s6y8473yrffqs9xlgntk0h
// g1wymu47drhr0kuq2098m792lytgtj2nyx77yrsm
