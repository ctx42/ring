// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"bytes"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"

	"github.com/ctx42/ring/internal/meta"
)

func Test_WithEnv(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithEnv([]string{"A=1", "B=2"})(rng)

	// --- Then ---
	assert.Equal(t, []string{"A=1", "B=2"}, Sort(rng.EnvGetAll()))
}

func Test_WithArgs(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithArgs([]string{"A=1", "B=2"})(rng)

	// --- Then ---
	assert.Equal(t, []string{"A=1", "B=2"}, rng.Args())
}

func Test_WithMeta(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithMeta(map[string]any{"A": 1, "B": 2})(rng)

	// --- Then ---
	assert.Equal(t, map[string]any{"A": 1, "B": 2}, rng.MetaGetAll())
}

func Test_WithClock(t *testing.T) {
	// --- Given ---
	rng := &Ring{}

	// --- When ---
	WithClock(time.Now)(rng)

	// --- Then ---
	assert.Same(t, time.Now, rng.Clock())
}

func Test_defaultRing(t *testing.T) {
	// --- When ---
	have := defaultRing()

	// --- Then ---
	assert.Nil(t, have.EnvGetAll())
	assert.Nil(t, have.MetaGetAll())
	assert.Same(t, os.Stdin, have.Stdin())
	assert.Same(t, os.Stdout, have.Stdout())
	assert.Same(t, os.Stderr, have.Stderr())
	assert.Within(t, time.Now(), "1ms", have.Clock()())
	assert.Zone(t, time.UTC, have.Clock()().Location())
	assert.Equal(t, os.Args[1:], have.Args())
	assert.Equal(t, os.Args[0], have.Name())

	// If this fails, that means the above assertions need to be adjusted.
	assert.Equal(t, 6, reflect.TypeOf(have).NumField())
}

func Test_New(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		// --- When ---
		have := New()

		// --- Then ---
		want := os.Environ()
		sort.Strings(want)
		assert.Equal(t, want, Sort(have.EnvGetAll()))
		assert.NotNil(t, have.MetaGetAll())
		assert.Same(t, os.Stdin, have.Stdin())
		assert.Same(t, os.Stdout, have.Stdout())
		assert.Same(t, os.Stderr, have.Stderr())
		assert.Within(t, time.Now(), "1ms", have.Clock()())
		assert.Zone(t, time.UTC, have.Clock()().Location())
		assert.Equal(t, os.Args[1:], have.Args())
		assert.Equal(t, os.Args[0], have.Name())

		// If this fails, that means the above assertions need to be adjusted.
		assert.Equal(t, 6, reflect.TypeOf(have).NumField())
	})

	t.Run("environment not overwritten", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2"}

		// --- When ---
		rng := New(WithEnv(env))

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(rng.EnvGetAll()))
	})
}

func Test_Ring_WithStdin_Stdin(t *testing.T) {
	// --- Given ---
	buf := &bytes.Buffer{}
	rng := New()

	// --- When ---
	have := rng.WithStdin(buf)

	// --- Then ---
	assert.Same(t, buf, have.Stdin())
	assert.NotSame(t, rng.Stdin(), have.Stdin())
}

func Test_Ring_WithStdout_Stdout(t *testing.T) {
	// --- Given ---
	buf := &bytes.Buffer{}
	rng := New()

	// --- When ---
	have := rng.WithStdout(buf)

	// --- Then ---
	assert.Same(t, buf, have.Stdout())
	assert.NotSame(t, rng.Stdout(), have.Stdout())
}

func Test_Ring_WithStderr_Stderr(t *testing.T) {
	// --- Given ---
	buf := &bytes.Buffer{}
	rng := New()

	// --- When ---
	have := rng.WithStderr(buf)

	// --- Then ---
	assert.Same(t, buf, have.Stderr())
	assert.NotSame(t, rng.Stderr(), have.Stderr())
}

func Test_Ring_Clock(t *testing.T) {
	// --- Given ---
	rng := New()

	// --- When ---
	have := rng.Clock()

	// --- Then ---
	assert.Within(t, time.Now(), "1ms", have())
	assert.Zone(t, time.UTC, have().Location())
}

func Test_Ring_WithArgs_Args(t *testing.T) {
	t.Run("set args", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2"}
		rng := New()

		// --- When ---
		have := rng.WithArgs(env)

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2"}, have.Args())
		assert.NotEqual(t, rng.Args(), have.Args())
	})

	t.Run("WithArgs returns clone", func(t *testing.T) {
		// --- Given ---
		rng0 := New(WithArgs([]string{}))

		// --- When ---
		rng1 := rng0.WithArgs([]string{"A=1", "B=2"})
		rng2 := rng1.WithArgs([]string{"C=3", "D=4"})
		rng3 := rng2.WithName("abc")

		// --- Then ---
		assert.Len(t, 0, rng0.Args())
		assert.Equal(t, []string{"A=1", "B=2"}, rng1.Args())
		assert.Equal(t, []string{"C=3", "D=4"}, rng2.Args())
		assert.Equal(t, []string{"C=3", "D=4"}, rng3.Args())
	})
}

func Test_Ring_WithName_Name(t *testing.T) {
	// --- Given ---
	rng := New()

	// --- When ---
	have := rng.WithName("my")

	// --- Then ---
	assert.Equal(t, "my", have.Name())
	assert.NotEqual(t, rng.Name(), have.Name())
}

func Test_Ring_Env(t *testing.T) {
	t.Run("get all", func(t *testing.T) {
		// --- Given ---
		rng := New(WithEnv([]string{"A=1", "B=2"}))

		// --- When ---
		env := rng.Env()

		// --- Then ---
		want := []string{"A=1", "B=2"}
		assert.Equal(t, want, Sort(env.EnvGetAll()))
	})

	t.Run("gets copy", func(t *testing.T) {
		// --- Given ---
		rng := New(WithEnv([]string{"A=1", "B=2"}))
		env := rng.Env()

		// --- When ---
		env.EnvSet("C", "3")

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2", "C=3"}, Sort(env.EnvGetAll()))
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(rng.EnvGetAll()))
	})
}

func Test_Ring_EnvSet(t *testing.T) {
	t.Run("set env", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2"}
		rng := New(WithEnv(env))

		// --- When ---
		have := rng.EnvSet("C", "3")

		// --- Then ---
		want := []string{"A=1", "B=2", "C=3"}
		assert.Equal(t, want, Sort(have.EnvGetAll()))
		assert.NotEqual(t, rng.EnvGetAll(), have.EnvGetAll())
	})

	t.Run("WithEnv returns clone", func(t *testing.T) {
		// --- Given ---
		rng0 := New(WithEnv([]string{}))

		// --- When ---
		rng1 := rng0.EnvSet("A", "1")
		rng2 := rng1.EnvSet("B", "2")
		rng3 := rng2.WithName("abc")

		// --- Then ---
		assert.Len(t, 0, rng0.EnvGetAll())
		assert.Equal(t, []string{"A=1"}, rng1.EnvGetAll())
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(rng2.EnvGetAll()))
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(rng3.EnvGetAll()))
	})
}

func Test_Ring_EnvSetBulk(t *testing.T) {
	// --- Given ---
	env := []string{"A=1", "B=2"}
	rng := New(WithEnv(env))

	// --- When ---
	have := rng.EnvSetBulk([]string{"A=A", "C=3"})

	// --- Then ---
	want := []string{"A=A", "B=2", "C=3"}
	assert.Equal(t, want, Sort(have.EnvGetAll()))
	assert.NotEqual(t, rng.EnvGetAll(), have.EnvGetAll())
}

func Test_Ring_EnvUnset(t *testing.T) {
	t.Run("unset env", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2"}
		rng := New(WithEnv(env))

		// --- When ---
		have := rng.EnvUnset("A")

		// --- Then ---
		assert.Equal(t, []string{"B=2"}, have.EnvGetAll())
		assert.NotEqual(t, rng.EnvGetAll(), have.EnvGetAll())
	})

	t.Run("EnvUnset returns clone", func(t *testing.T) {
		// --- Given ---
		rng0 := New(WithEnv([]string{"A=1", "B=2"}))

		// --- When ---
		rng1 := rng0.EnvUnset("A")
		rng2 := rng1.EnvUnset("B")
		rng3 := rng2.WithName("abc")

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(rng0.EnvGetAll()))
		assert.Equal(t, []string{"B=2"}, rng1.EnvGetAll())
		assert.Nil(t, rng2.EnvGetAll())
		assert.Nil(t, rng3.EnvGetAll())
	})
}

func Test_Ring_MetaSet(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		rng := New()

		// --- When ---
		have := rng.MetaSet("A", 1)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, rng.MetaGetAll())
		assert.Equal(t, map[string]any{"A": 1}, have.MetaGetAll())
	})

	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1}))

		// --- When ---
		have := rng.MetaSet("A", 2)

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 2}, rng.MetaGetAll())
		assert.Equal(t, map[string]any{"A": 2}, have.MetaGetAll())
	})
}

func Test_Ring_MetaDelete(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1, "B": 2}))

		// --- When ---
		have := rng.MetaDelete("A")

		// --- Then ---
		assert.Equal(t, map[string]any{"B": 2}, have.MetaGetAll())
		assert.Equal(t, map[string]any{"B": 2}, rng.MetaGetAll())
	})

	t.Run("delete not existing", func(t *testing.T) {
		// --- Given ---
		rng := New(WithMeta(map[string]any{"A": 1}))

		// --- When ---
		have := rng.MetaDelete("B")

		// --- Then ---
		assert.Equal(t, map[string]any{"A": 1}, have.MetaGetAll())
		assert.Equal(t, map[string]any{"A": 1}, rng.MetaGetAll())
		assert.Same(t, rng.MetaGetAll(), have.MetaGetAll())
	})
}

func Test_Ring_IsZero(t *testing.T) {
	t.Run("is zero", func(t *testing.T) {
		// --- Given ---
		rng := Ring{}

		// --- When ---
		have := rng.IsZero()

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("new is not", func(t *testing.T) {
		// --- Given ---
		rng := New()

		// --- When ---
		have := rng.IsZero()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not zero if env set", func(t *testing.T) {
		// --- Given ---
		rng := Ring{hidEnv: NewEnv(nil)}

		// --- When ---
		have := rng.IsZero()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not zero if meta set", func(t *testing.T) {
		// --- Given ---
		rng := Ring{hidMeta: meta.New()}

		// --- When ---
		have := rng.IsZero()

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("not zero if name set", func(t *testing.T) {
		// --- Given ---
		rng := Ring{name: "abc"}

		// --- When ---
		have := rng.IsZero()

		// --- Then ---
		assert.False(t, have)
	})
}
