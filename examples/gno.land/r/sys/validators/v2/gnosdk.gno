package validators

import (
	"math"

	"gno.land/p/sys/validators"
)

// GetChanges returns the validator changes stored on the realm,
// for blocks in the [from, to] range (inclusive on both ends).
// Panics if from > to or if to >= math.MaxInt64 (to avoid overflow).
// This function is intended to be called by gno.land through the GnoSDK.
func GetChanges(from, to int64) []validators.Validator {
	if to > math.MaxInt64-1 {
		panic("to exceeds maximum allowed value")
	}
	if to < from {
		panic("invalid range: from must be <= to")
	}

	valsetChanges := make([]validators.Validator, 0)

	// Gather the changes in the [from, to] block range.
	// AVL Iterate uses an exclusive end, so we pass to+1.
	changes.Iterate(getBlockID(from), getBlockID(to+1), func(_ string, value any) bool {
		chs := value.([]change)

		for _, ch := range chs {
			valsetChanges = append(valsetChanges, ch.validator)
		}

		return false
	})

	return valsetChanges
}
