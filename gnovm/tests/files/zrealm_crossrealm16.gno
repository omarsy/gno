// PKGPATH: gno.land/r/crossrealm
package crossrealm

import (
	"std"

	crossrealm "gno.land/r/demo/tests/crossrealm"
)

type fooer struct{}

func (fooer) Foo(cur realm) {
	println("hello " + std.CurrentRealm().PkgPath())
}

func (fooer) Bar() {
	println("hello " + std.CurrentRealm().PkgPath())
}

var f *fooer

func main(cur realm) {
	f = &fooer{}
	crossrealm.SetFooer(cross, f)
	// crossrealm.CallFooerFooCur(cross) <-- panics: cannot cur-call to
	// external realm function gno.land/r/crossrealm_test.Foo from
	// gno.land/r/demo/tests/crossrealm
	crossrealm.CallFooerFooCross(cross)
	crossrealm.CallFooerBar()
	crossrealm.CallFooerBarCrossing(cross)
	print(".")
}

// Output:
// hello gno.land/r/crossrealm
// hello gno.land/r/crossrealm
// hello gno.land/r/demo/tests/crossrealm
// .
