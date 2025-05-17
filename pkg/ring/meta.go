// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"github.com/ctx42/ring/internal/meta"
)

// Metadata is an interface to implement by types which deal with metadata.
type Metadata interface {
	// MetaLookup returns the value of the variable named by the key. If the
	// variable is present in the map, the value (which may be empty or nil) is
	// returned and the boolean is true. Otherwise, the returned value will be
	// nil and the boolean will be false.
	MetaLookup(key string) (any, bool)

	// MetaGet returns the value of the variable named by the key. If the
	// variable is not present, in the map nil is returned.
	MetaGet(key string) any

	// MetaSet sets the value of variable named by the key.
	MetaSet(key string, value any)

	// MetaDelete deletes the map entry identified by the key.
	MetaDelete(key string)

	// MetaGetAll returns the underlying map used by [Meta]. After call to this
	// method [Meta] instance must no longer be used.
	MetaGetAll() map[string]any
}

var _ Metadata = meta.Meta{} // Compile time check.
