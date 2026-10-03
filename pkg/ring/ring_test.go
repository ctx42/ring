// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_WithEnv(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	env := []string{"A=1", "B=2"}

	// --- When ---
	WithEnv(env)(rng)

	// --- Then ---
	assert.Equal(t, map[string]string{"A": "1", "B": "2"}, rng.hidEnv.env)

	t.Run("nil slice", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{}

		var env []string

		// --- When ---
		WithEnv(env)(rng)

		// --- Then ---
		assert.NotNil(t, rng.hidEnv)
		assert.NotNil(t, rng.hidEnv.env)
		assert.Empty(t, rng.hidEnv.env)
	})
}

func Test_WithName(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	name := "abc"

	// --- When ---
	WithName(name)(rng)

	// --- Then ---
	assert.Equal(t, name, rng.name)
}

func Test_WithArgs(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	args := []string{"A=1", "B=2"}

	// --- When ---
	WithArgs(args)(rng)

	// --- Then ---
	assert.Equal(t, args, rng.args)
}

func Test_WithClock(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	clk := time.Now

	// --- When ---
	WithClock(clk)(rng)

	// --- Then ---
	assert.Same(t, clk, rng.clock)
}

func Test_WithMeta(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	meta := map[string]any{"A": 1, "B": 2}

	// --- When ---
	WithMeta(meta)(rng)

	// --- Then ---
	assert.Equal(t, meta, rng.meta)
}

func Test_WithFS(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	root := must.Value(os.OpenRoot("ringtest"))
	t.Cleanup(func() { _ = root.Close() })

	FS := root.FS()

	// --- When ---
	WithFS(FS)(rng)

	// --- Then ---
	assert.Same(t, FS, rng.fs)
}

func Test_defaultRing(t *testing.T) {
	t.Run("empty os.Args", func(t *testing.T) {
		// --- Given ---
		orig := os.Args
		t.Cleanup(func() { os.Args = orig })
		os.Args = nil

		// --- When ---
		have := defaultRing()

		// --- Then ---
		assert.Equal(t, "", have.name)
		assert.Nil(t, have.args)
	})

	// --- When ---
	have := defaultRing()

	// --- Then ---
	assert.Nil(t, have.hidEnv)
	assert.Same(t, os.Stdin, have.stdin)
	assert.Same(t, os.Stdout, have.stdout)
	assert.Same(t, os.Stderr, have.stderr)
	assert.Same(t, NowUTC, have.clock)
	assert.Nil(t, have.fs)
	assert.Equal(t, os.Args[0], have.name)
	assert.Equal(t, os.Args[1:], have.args)
	assert.Nil(t, have.meta)

	assert.Fields(t, 7, Ring{})
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
		assert.Nil(t, have.fs)
		assert.Equal(t, os.Args[0], have.name)
		assert.Equal(t, os.Args[1:], have.args)
		assert.NotNil(t, have.meta)
		assert.Empty(t, have.meta)

		assert.Fields(t, 7, Ring{})
	})

	t.Run("with option", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2"}

		opt := WithEnv(env)

		// --- When ---
		have := New(opt)

		// --- Then ---
		assert.Equal(t, map[string]string{"A": "1", "B": "2"}, have.env)
	})

	t.Run("nil option", func(t *testing.T) {
		// --- Given ---
		name := "app"

		opt := WithName(name)

		opts := []Option{nil, opt}

		// --- When ---
		have := New(opts...)

		// --- Then ---
		assert.Equal(t, name, have.Name())
	})

	t.Run("nil clock", func(t *testing.T) {
		// --- Given ---
		var clk Clock

		opt := WithClock(clk)

		// --- When ---
		have := New(opt)

		// --- Then ---
		assert.Same(t, NowUTC, have.Clock())
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

		key := "A"

		val := 1

		// --- When ---
		rng.MetaSet(key, val)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, rng.meta)
	})

	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1}))

		key := "A"

		val := 2

		// --- When ---
		rng.MetaSet(key, val)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 2}, rng.meta)
	})

	t.Run("nil map", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{}

		key := "A"

		val := 1

		// --- When ---
		rng.MetaSet(key, val)

		// --- Then ---
		assert.Equal(t, val, rng.MetaGet(key))
	})
}

func Test_Ring_MetaGet(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{"A": 1}}

		key := "A"

		// --- When ---
		have := rng.MetaGet(key)

		// --- Then ---
		assert.Equal(t, 1, have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{}}

		key := "B"

		// --- When ---
		have := rng.MetaGet(key)

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_Ring_MetaLookup(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{"A": 1}}

		key := "A"

		// --- When ---
		have, ok := rng.MetaLookup(key)

		// --- Then ---
		assert.Equal(t, 1, have)
		assert.True(t, ok)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{meta: map[string]any{}}

		key := "B"

		// --- When ---
		have, ok := rng.MetaLookup(key)

		// --- Then ---
		assert.Nil(t, have)
		assert.False(t, ok)
	})
}

func Test_Ring_MetaDelete(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1, "B": 2}))

		key := "A"

		// --- When ---
		rng.MetaDelete(key)

		// --- Then ---
		assert.Equal(t, map[string]any{"B": 2}, rng.meta)
	})

	t.Run("delete not existing", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1}))

		key := "B"

		// --- When ---
		rng.MetaDelete(key)

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

func Test_Ring_FS(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{fs: os.DirFS("ringtest")}

		// --- When ---
		have, err := rng.FS()

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, rng.fs, have)
	})

	t.Run("error - no filesystem access", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{}

		// --- When ---
		have, err := rng.FS()

		// --- Then ---
		assert.ErrorIs(t, ErrNoFsAccess, err)
		assert.Nil(t, have)
	})
}

func Test_Ring_Clone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		rngFS := os.DirFS("ringtest")

		sin := &bytes.Buffer{}
		sout := &bytes.Buffer{}
		eout := &bytes.Buffer{}

		rng := New(WithFS(rngFS))
		rng.SetStdin(sin)
		rng.SetStdout(sout)
		rng.SetStderr(eout)

		// --- When ---
		have := rng.Clone()

		// --- Then ---
		assert.NotSame(t, rng, have)
		assert.NotSame(t, rng.hidEnv, have.hidEnv)
		assert.Equal(t, rng.hidEnv.env, have.hidEnv.env)
		assert.NotSame(t, rng.hidEnv.env, have.hidEnv.env)
		assert.NotSame(t, rng.hidIO, have.hidIO)
		assert.Same(t, rng.Stdin(), have.Stdin())
		assert.Same(t, rng.Stdout(), have.Stdout())
		assert.Same(t, rng.Stderr(), have.Stderr())
		assert.Same(t, rng.clock, have.clock)
		assert.Equal(t, rng.name, have.name)
		assert.Equal(t, rngFS, have.fs)
		assert.Equal(t, rng.args, have.args)
		assert.NotSame(t, rng.args, have.args)
		assert.Same(t, rng.meta, have.meta)

		assert.Fields(t, 7, Ring{})
	})

	t.Run("zero value", func(t *testing.T) {
		// --- Given ---
		rng := &Ring{}

		// --- When ---
		have := rng.Clone()

		// --- Then ---
		assert.Nil(t, have.hidEnv)
		assert.Nil(t, have.hidIO)
		assert.NotSame(t, rng, have)
	})
}
