// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"os"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_WithEnv(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithEnv([]string{"A=1", "B=2"})(rng)

	// --- Then ---
	assert.Equal(t, map[string]string{"A": "1", "B": "2"}, rng.hidEnv.env)
}

func Test_WithName(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithName("abc")(rng)

	// --- Then ---
	assert.Equal(t, "abc", rng.name)
}

func Test_WithArgs(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithArgs([]string{"A=1", "B=2"})(rng)

	// --- Then ---
	assert.Equal(t, []string{"A=1", "B=2"}, rng.args)
}

func Test_WithMeta(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithMeta(map[string]any{"A": 1, "B": 2})(rng)

	// --- Then ---
	assert.Equal(t, map[string]any{"A": 1, "B": 2}, rng.meta)
}

func Test_WithClock(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithClock(time.Now)(rng)

	// --- Then ---
	assert.Same(t, time.Now, rng.clock)
}

func Test_defaultRing(t *testing.T) {
	// --- When ---
	have := defaultRing()

	// --- Then ---
	assert.Nil(t, have.hidEnv)
	assert.Same(t, os.Stdin, have.stdin)
	assert.Same(t, os.Stdout, have.stdout)
	assert.Same(t, os.Stderr, have.stderr)
	assert.Same(t, NowUTC, have.clock)
	assert.Equal(t, os.Args[0], have.name)
	assert.Equal(t, os.Args[1:], have.args)
	assert.Nil(t, have.meta)
	assert.Fields(t, 6, Ring{})
}

func Test_New(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have := New()

		// --- Then ---
		assert.Equal(t, Sort(os.Environ()), Sort(have.EnvAll()))
		assert.Same(t, os.Stdin, have.stdin)
		assert.Same(t, os.Stdout, have.stdout)
		assert.Same(t, os.Stderr, have.stderr)
		assert.Same(t, NowUTC, have.clock)
		assert.Equal(t, os.Args[0], have.name)
		assert.Equal(t, os.Args[1:], have.args)
		assert.NotNil(t, have.meta)
		assert.Empty(t, have.meta)
		assert.Fields(t, 6, Ring{})
	})

	t.Run("with option", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2"}

		// --- When ---
		rng := New(WithEnv(env))

		// --- Then ---
		assert.Equal(t, map[string]string{"A": "1", "B": "2"}, rng.env)
	})
}

func Test_Ring_Clock(t *testing.T) {
	// --- Given ---
	custom := func() time.Time { return time.Time{} }
	rng := &Ring{clock: custom}

	// --- When ---
	have := rng.Clock()

	// --- Then ---
	assert.Same(t, custom, have)
}

func Test_Ring_Args(t *testing.T) {
	// --- Given ---
	args := []string{"-arg0", "-arg1"}
	rng := &Ring{args: args}

	// --- When ---
	have := rng.Args()

	// --- Then ---
	assert.Same(t, args, have)
}

func Test_Ring_SetArgs(t *testing.T) {
	// --- Given ---
	args := []string{"-arg0", "-arg1"}
	rng := &Ring{}

	// --- When ---
	have := rng.SetArgs(args)

	// --- Then ---
	assert.Same(t, rng, have)
	assert.Same(t, args, rng.args)
}

func Test_Ring_Name(t *testing.T) {
	// --- Given ---
	rng := &Ring{name: "abc"}

	// --- When ---
	have := rng.Name()

	// --- Then ---
	assert.Equal(t, "abc", have)
}

func Test_Ring_MetaSet(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		rng := New()

		// --- When ---
		rng.MetaSet("A", 1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, rng.meta)
	})

	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1}))

		// --- When ---
		rng.MetaSet("A", 2)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 2}, rng.meta)
	})
}

func Test_Ring_MetaGet(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{"A": 1}}

		// --- When ---
		have := rng.MetaGet("A")

		// --- Then ---
		assert.Equal(t, 1, have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{}}

		// --- When ---
		have := rng.MetaGet("B")

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_Ring_MetaLookup(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{"A": 1}}

		// --- When ---
		have, ok := rng.MetaLookup("A")

		// --- Then ---
		assert.Equal(t, 1, have)
		assert.True(t, ok)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{}}

		// --- When ---
		have, ok := rng.MetaLookup("B")

		// --- Then ---
		assert.Nil(t, have)
		assert.False(t, ok)
	})
}

func Test_Ring_MetaDelete(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1, "B": 2}))

		// --- When ---
		rng.MetaDelete("A")

		// --- Then ---
		assert.Equal(t, map[string]any{"B": 2}, rng.meta)
	})

	t.Run("delete not existing", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1}))

		// --- When ---
		rng.MetaDelete("B")

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, rng.meta)
	})
}

func Test_Ring_MetaAll(t *testing.T) {
	// --- Given ---
	rng := &Ring{meta: map[string]any{"A": 1}}

	// --- When ---
	have := rng.MetaAll()

	// --- Then ---
	assert.Equal(t, map[string]any{"A": 1}, have)
	assert.Same(t, rng.meta, have)
}
