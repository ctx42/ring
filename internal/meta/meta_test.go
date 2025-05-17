// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package meta

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithLen(t *testing.T) {
	// --- Given ---
	opts := &metaOpts{}

	// --- When ---
	WithLen(10)(opts)

	// --- Then ---
	assert.Equal(t, 10, opts.length)
}

func Test_WithMap(t *testing.T) {
	// --- Given ---
	want := map[string]any{"A": 1}
	opts := &metaOpts{}

	// --- When ---
	WithMap(want)(opts)

	// --- Then ---
	assert.Same(t, want, opts.initial)
}

func Test_New(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have := New()

		// --- Then ---
		assert.NotNil(t, have.m)
		assert.Len(t, 0, have.m)
	})

	t.Run("with initial map", func(t *testing.T) {
		// --- Given ---
		want := map[string]any{"A": 1}

		// --- When ---
		have := New(WithMap(want))

		// --- Then ---
		assert.Same(t, want, have.m)
	})
}

func Test_Meta_MetaSet(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		m := New()

		// --- When ---
		m.MetaSet("A", 1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, m.m)
	})

	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1}))

		// --- When ---
		m.MetaSet("A", 2)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 2}, m.m)
	})
}

func Test_Meta_MetaLookup(t *testing.T) {
	t.Run("empty collection", func(t *testing.T) {
		// --- Given ---
		m := New()

		// --- When ---
		haveVal, haveExi := m.MetaLookup("A")

		// --- Then ---
		assert.False(t, haveExi)
		assert.Nil(t, haveVal)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1}))

		// --- When ---
		haveVal, haveExi := m.MetaLookup("A")

		// --- Then ---
		assert.True(t, haveExi)
		assert.Equal(t, 1, haveVal)
	})
}

func Test_Meta_MetaGet(t *testing.T) {
	t.Run("empty collection", func(t *testing.T) {
		// --- Given ---
		m := New()

		// --- When ---
		have := m.MetaGet("A")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1}))

		// --- When ---
		have := m.MetaGet("A")

		// --- Then ---
		assert.Equal(t, 1, have)
	})
}

func Test_Meta_MetaDelete(t *testing.T) {
	t.Run("delete not existing", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1, "B": 2, "C": 3}))

		// --- When ---
		m.MetaDelete("D")

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1, "B": 2, "C": 3}, m.MetaGetAll())
	})

	t.Run("delete existing", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1, "B": 2, "C": 3}))

		// --- When ---
		m.MetaDelete("A")

		// --- Then ---
		assert.Equal(t, map[string]any{"B": 2, "C": 3}, m.MetaGetAll())
	})
}

func Test_Meta_MetaGetAll(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		// --- Given ---
		m := New()

		// --- When ---
		have := m.MetaGetAll()

		// --- Then ---
		assert.NotNil(t, have)
		assert.Len(t, 0, have)
	})

	t.Run("not empty", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1, "B": 2}))

		// --- When ---
		have := m.MetaGetAll()

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1, "B": 2}, have)
	})
}

func Test_Meta_MetaLen(t *testing.T) {
	t.Run("empty collection", func(t *testing.T) {
		// --- Given ---
		m := New()

		// --- When ---
		have := m.MetaLen()

		// --- Then ---
		assert.Equal(t, 0, have)
	})

	t.Run("collection with keys", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1, "B": 2}))

		// --- When ---
		have := m.MetaLen()

		// --- Then ---
		assert.Equal(t, 2, have)
	})
}

func Test_Meta_MetaIsNil(t *testing.T) {
	t.Run("declared", func(t *testing.T) {
		// --- Given ---
		var m Meta

		// --- When ---
		have := m.MetaIsNil()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("constructed", func(t *testing.T) {
		// --- Given ---
		m := New()

		// --- When ---
		have := m.MetaIsNil()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("collection with keys", func(t *testing.T) {
		// --- Given ---
		m := New(WithMap(map[string]any{"A": 1, "B": 2}))

		// --- When ---
		have := m.MetaIsNil()

		// --- Then ---
		assert.False(t, have)
	})
}
