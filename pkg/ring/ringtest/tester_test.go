// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ringtest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_WithEnv(t *testing.T) {
	// --- Given ---
	env := []string{"A=B", "C=D"}
	tst := &Tester{}

	// --- When ---
	WithEnv(env)(tst)

	// --- Then ---
	assert.Equal(t, env, Sort(tst.EnvGetAll()))
}

func Test_WithName(t *testing.T) {
	// --- Given ---
	tst := &Tester{}

	// --- When ---
	WithName("my")(tst)

	// --- Then ---
	assert.Equal(t, "my", tst.Name())
}

func Test_WithMeta(t *testing.T) {
	// --- Given ---
	m := map[string]any{"A": 1}
	tst := &Tester{}

	// --- When ---
	WithMeta(m)(tst)

	// --- Then ---
	assert.Same(t, m, tst.m)
}

func Test_WithClock(t *testing.T) {
	// --- Given ---
	fn := func() time.Time { return time.Now().UTC() }
	tst := &Tester{}

	// --- When ---
	WithClock(fn)(tst)

	// --- Then ---
	assert.Same(t, fn, tst.Clock())
}

func Test_New(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		// --- When ---
		tst := New(tspy)

		// --- Then ---
		want := os.Environ()
		sort.Strings(want)
		assert.Equal(t, want, Sort(tst.EnvGetAll()))
		assert.NotNil(t, tst.m)
		assert.Empty(t, tst.m)
		content := must.Value(io.ReadAll(tst.Stdin()))
		assert.Equal(t, "", string(content))
		assert.Equal(t, "", tst.Stdout())
		assert.Equal(t, "", tst.Stderr())
		assert.Within(t, time.Now(), "1ms", tst.Clock()())
		assert.Zone(t, time.UTC, tst.Clock()().Location())
		assert.Equal(t, os.Args[0], tst.Name())
		assert.Same(t, tspy, tst.t)
	})

	t.Run("with environment", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		env := []string{"A=B", "C=D"}

		// --- When ---
		tst := New(tspy, WithEnv(env))

		// --- Then ---
		assert.Equal(t, env, Sort(tst.EnvGetAll()))
	})
}

func Test_Tester_Ring(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		rng := tst.Ring("a", "b", "c")

		// --- Then ---
		assert.Equal(t, Sort(tst.EnvGetAll()), Sort(rng.EnvGetAll()))
		assert.NotSame(t, tst.EnvGetAll(), rng.EnvGetAll())
		assert.Equal(t, tst.m, rng.MetaAll())
		assert.NotSame(t, tst.m, rng.MetaAll())
		assert.Same(t, tst.sin, rng.Stdin())
		assert.Same(t, tst.sout, rng.Stdout())
		assert.Same(t, tst.eout, rng.Stderr())
		assert.Same(t, tst.clock, rng.Clock())
		assert.Equal(t, os.Args[0], rng.Name())
		assert.Equal(t, []string{"a", "b", "c"}, rng.Args())
	})

	t.Run("with custom name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy, WithName("my"))

		// --- When ---
		rng := tst.Ring("a", "b", "c")

		// --- Then ---
		assert.Equal(t, "my", rng.Name())
	})
}

func Test_Tester_Streams(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)

	// --- When ---
	ios := tst.Streams()

	// --- Then ---
	assert.Same(t, tst.sin, ios.Stdin())
	assert.Same(t, tst.sout, ios.Stdout())
	assert.Same(t, tst.eout, ios.Stderr())
}

func Test_Tester_WithStdin(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	buf := bytes.NewBuffer([]byte("test"))
	tst := New(tspy)

	// --- When ---
	have := tst.WithStdin(buf)

	// --- Then ---
	assert.Same(t, tst, have)
	content := must.Value(io.ReadAll(tst.Stdin()))
	assert.Equal(t, "test", string(content))
}

func Test_Tester_WetStdout(t *testing.T) {
	t.Run("want stdout wet but is dry", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(3)
		tspy.ExpectError()
		tspy.ExpectLogEqual("expected buffer not to be empty:\n  name: stdout")
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		have := tst.WetStdout()

		// --- Then ---
		assert.Same(t, tst, have)
	})

	t.Run("want stdout dry but is wet", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.ExpectError()
		wMsg := "expected buffer to be empty:\n" +
			"  name: stdout\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		_, _ = tst.sout.WriteString("abc")
	})
}

func Test_Tester_ResetStdout(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)
	_, _ = fmt.Fprint(tst.sout, "test")

	// --- When ---
	tst.ResetStdout()
}

func Test_Tester_WetStderr(t *testing.T) {
	t.Run("want stderr wet but is dry", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(3)
		tspy.ExpectError()
		tspy.ExpectLogEqual("expected buffer not to be empty:\n  name: stderr")
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		have := tst.WetStderr()

		// --- Then ---
		assert.Same(t, tst, have)
	})

	t.Run("want stderr dry but is wet", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.ExpectError()
		wMsg := "expected buffer to be empty:\n" +
			"  name: stderr\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		_, _ = tst.eout.WriteString("abc")
	})
}

func Test_Tester_ResetStderr(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)
	_, _ = fmt.Fprint(tst.eout, "test")

	// --- When ---
	tst.ResetStderr()
}

func Test_Tester_EnvGetAll(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy, WithEnv([]string{"A=B", "C=D"}))

	// --- When ---
	have := Sort(tst.EnvGetAll())

	// --- Then ---
	assert.Equal(t, []string{"A=B", "C=D"}, have)
}

func Test_Tester_Name(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy, WithName("name"))

	// --- When ---
	have := tst.Name()

	// --- Then ---
	assert.Equal(t, "name", have)
}

func Test_Tester_Stdout(t *testing.T) {
	t.Run("nothing written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		have := tst.Stdout()

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("data written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(3)
		tspy.Close()

		tst := New(tspy).WetStdout()
		_, _ = tst.sout.WriteString("abc")

		// --- When ---
		have := tst.Stdout()

		// --- Then ---
		assert.Equal(t, "abc", have)
	})
}

func Test_Tester_Stderr(t *testing.T) {
	t.Run("nothing written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		have := tst.Stderr()

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("data written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(3)
		tspy.Close()

		tst := New(tspy).WetStderr()
		_, _ = tst.eout.WriteString("abc")

		// --- When ---
		have := tst.Stderr()

		// --- Then ---
		assert.Equal(t, "abc", have)
	})
}

func Test_Tester_Clock(t *testing.T) {
	t.Run("nothing written", func(t *testing.T) {
		// --- Given ---
		tim := time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
		clk := func() time.Time { return tim }

		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy, WithClock(clk))

		// --- When ---
		have := tst.Clock()

		// --- Then ---
		assert.Exact(t, tim, have())
	})
}
