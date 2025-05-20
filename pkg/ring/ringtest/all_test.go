package ringtest

import (
	"sort"
)

// Sort works like [sort.Strings] but returns the sorted slice.
func Sort(in []string) []string {
	sort.Strings(in)
	return in
}
