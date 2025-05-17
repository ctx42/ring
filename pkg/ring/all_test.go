// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"sort"
)

// Sort works like [sort.Strings] but returns sorted slice.
func Sort(in []string) []string {
	sort.Strings(in)
	return in
}
