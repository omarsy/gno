package validators

import (
	"gno.land/p/sys/validators"
)

// GetChanges returns the validator changes stored on the realm,
// for blocks in the [from, to] range (inclusive on both ends).
// This function is intended to be called by gno.land through the GnoSDK.
func GetChanges(from, to int64) []validators.Validator {
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
