// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package meta

// WithLen is option for [New] setting default length for the meta map.
func WithLen(n int) func(*metaOpts) {
	return func(o *metaOpts) { o.length = n }
}

// WithMap is an option for [New] setting map to use for the metadata
// collection.
func WithMap(m map[string]any) func(opts *metaOpts) {
	return func(o *metaOpts) { o.initial = m }
}

// metaOpts represents [Meta] options used when creating the instance.
type metaOpts struct {
	length  int            // Initial metadata map size, default is 10.
	initial map[string]any // Initial metadata map.
}

// Meta represents metadata.
type Meta struct {
	m map[string]any
}

// New returns new [Meta] instance. By default, the new map is initialized with
// length equal to 10.
func New(opts ...func(*metaOpts)) Meta {
	def := &metaOpts{length: 10}
	for _, opt := range opts {
		opt(def)
	}
	m := Meta{m: def.initial}
	if m.m == nil {
		m.m = make(map[string]any, def.length)
	}
	return m
}

// MetaSet sets the value of variable named by the key.
func (m Meta) MetaSet(key string, value any) {
	m.m[key] = value
}

// MetaLookup returns the value of the variable named by the key. If the
// variable is present in the map, the value (which may be empty or nil) is
// returned and the boolean is true. Otherwise, the returned value will be nil
// and the boolean will be false.
func (m Meta) MetaLookup(key string) (any, bool) {
	val, ok := m.m[key]
	return val, ok
}

// MetaGet returns the value of the variable named by the key. If the variable
// is not present, in the map nil is returned.
func (m Meta) MetaGet(key string) any { return m.m[key] }

// MetaDelete deletes the map entry identified by the key.
func (m Meta) MetaDelete(key string) { delete(m.m, key) }

// MetaGetAll returns the underlying map used by [Meta]. After call to this
// method [Meta] instance must no longer be used.
func (m Meta) MetaGetAll() map[string]any { return m.m }

// MetaLen returns number of entries in the map.
func (m Meta) MetaLen() int { return len(m.m) }

// MetaIsNil returns true if no memory was allocated for the underlying map.
func (m Meta) MetaIsNil() bool { return m.m == nil }
