package foo721

import (
	"testing"

	"gno.land/p/demo/grc/grc721"
	"gno.land/r/demo/users"

	pusers "gno.land/p/demo/users"
)

func TestFoo721(t *testing.T) {
	admin := pusers.AddressOrName("g1us8428u2a5satrlxzagqqa5m6vmuze025anjlj")
	hariom := pusers.AddressOrName("g1var589z07ppjsjd24ukm4uguzwdt0tw7g47cgm")

	for _, tc := range []struct {
		name     string
		expected any
		fn       func() any
	}{
		{"BalanceOf(admin)", uint64(10), func() any { return BalanceOf(admin) }},
		{"BalanceOf(hariom)", uint64(5), func() any { return BalanceOf(hariom) }},
		{"OwnerOf(0)", users.Resolve(admin), func() any { return OwnerOf(grc721.TokenID("0")) }},
		{"IsApprovedForAll(admin, hariom)", false, func() any { return IsApprovedForAll(admin, hariom) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn()
			if tc.expected != got {
				t.Errorf("expected: %v got: %v", tc.expected, got)
			}
		})
	}
}
