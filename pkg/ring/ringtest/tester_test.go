// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ringtest

import (
	"bytes"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/tester"

	"github.com/ctx42/ring/pkg/ring"
)

func Test_New(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		// --- When ---
		have := New(tspy)

		// --- Then ---

		// The instance of [ring.Ring].
		assert.Equal(t, Sort(os.Environ()), Sort(have.rng.EnvAll()))
		assert.NotNil(t, have.rng.MetaAll())
		assert.Empty(t, have.rng.MetaAll())
		assert.Same(t, os.Stdin, have.rng.Stdin())
		assert.Same(t, os.Stdout, have.rng.Stdout())
		assert.Same(t, os.Stderr, have.rng.Stderr())
		assert.Same(t, ring.NowUTC, have.rng.Clock())
		assert.Equal(t, os.Args[0], have.rng.Name())
		assert.Nil(t, have.rng.Args())

		// The instance of [tester.Tester].
		assert.Empty(t, have.sin.String())
		assert.Equal(t, "", have.sout.String())
		assert.Equal(t, "", have.eout.String())
		assert.Same(t, tspy, have.t)
	})

	t.Run("with environment", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		env := []string{"A=B", "C=D"}

		opt := ring.WithEnv(env)

		// --- When ---
		have := New(tspy, opt)

		// --- Then ---
		assert.Equal(t, env, Sort(have.Ring().EnvAll()))
	})

	t.Run("error - stdout written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.ExpectError()
		wMsg := "" +
			"expected buffer to be empty:\n" +
			"  name: stdout\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		tst := New(tspy)

		text := "abc"

		// --- When ---
		_, _ = tst.sout.WriteString(text)

		// --- Then ---
		assert.Equal(t, "abc", tst.sout.String())
	})

	t.Run("error - stderr written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.ExpectError()
		wMsg := "" +
			"expected buffer to be empty:\n" +
			"  name: stderr\n" +
			"  want: <empty>\n" +
			"  have: abc"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		tst := New(tspy)

		text := "abc"

		// --- When ---
		_, _ = tst.eout.WriteString(text)

		// --- Then ---
		assert.Equal(t, "abc", tst.eout.String())
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
		have := tst.Ring()

		// --- Then ---
		assert.Equal(t, Sort(os.Environ()), Sort(have.EnvAll()))
		assert.NotNil(t, have.MetaAll())
		assert.Empty(t, have.MetaAll())
		assert.Same(t, tst.sin, have.Stdin())
		assert.Same(t, tst.sout, have.Stdout())
		assert.Same(t, tst.eout, have.Stderr())
		assert.Same(t, ring.NowUTC, have.Clock())
		assert.Equal(t, os.Args[0], have.Name())
		assert.Nil(t, have.Args())
	})

	t.Run("with args", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy)

		args := []string{"a", "b", "c"}

		// --- When ---
		have := tst.Ring(args...)

		// --- Then ---
		assert.Equal(t, Sort(os.Environ()), Sort(have.EnvAll()))
		assert.NotNil(t, have.MetaAll())
		assert.Empty(t, have.MetaAll())
		assert.Same(t, tst.sin, have.Stdin())
		assert.Same(t, tst.sout, have.Stdout())
		assert.Same(t, tst.eout, have.Stderr())
		assert.Same(t, ring.NowUTC, have.Clock())
		assert.Equal(t, os.Args[0], have.Name())
		assert.Equal(t, []string{"a", "b", "c"}, have.Args())
	})

	t.Run("with a custom name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy, ring.WithName("my"))

		// --- When ---
		have := tst.Ring()

		// --- Then ---
		assert.Equal(t, "my", have.Name())
	})

	t.Run("with clone of metadata", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		m := map[string]any{"A": 1}
		tst := New(tspy, ring.WithMeta(m))

		// --- When ---
		have := tst.Ring()

		// --- Then ---
		assert.NotSame(t, m, have.MetaAll())
	})

	t.Run("with filesystem", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		root := os.DirFS("..")

		tst := New(tspy, ring.WithFS(root))

		// --- When ---
		have := tst.Ring()

		// --- Then ---
		filesystem := must.Value(have.FS())
		assert.Equal(t, root, filesystem)
	})
}

func Test_Tester_Streams(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)

	// --- When ---
	have := tst.Streams()

	// --- Then ---
	assert.Same(t, tst.sin, have.Stdin())
	assert.Same(t, tst.sout, have.Stdout())
	assert.Same(t, tst.eout, have.Stderr())
}

func Test_Tester_SetStdin(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	buf := &bytes.Buffer{}

	tst := New(tspy)

	// --- When ---
	have := tst.SetStdin(buf)

	// --- Then ---
	assert.Same(t, tst, have)

	assert.Same(t, buf, tst.sin)
}

func Test_Tester_WetStdout(t *testing.T) {
	t.Run("error - stdout dry", func(t *testing.T) {
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
}

func Test_Tester_ResetStdout(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)
	_, _ = tst.sout.WriteString("test")

	// --- When ---
	tst.ResetStdout()

	// --- Then ---
	assert.Equal(t, "", tst.sout.String())
}

func Test_Tester_WetStderr(t *testing.T) {
	t.Run("error - stderr dry", func(t *testing.T) {
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
}

func Test_Tester_ResetStderr(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)
	_, _ = tst.eout.WriteString("test")

	// --- When ---
	tst.ResetStderr()

	// --- Then ---
	assert.Equal(t, "", tst.eout.String())
}

func Test_Tester_Stdin(t *testing.T) {
	// --- Given ---
	tspy := tester.New(t)
	tspy.ExpectCleanups(2)
	tspy.Close()

	tst := New(tspy)
	_, _ = tst.sin.WriteString("abc")

	// --- When ---
	have := tst.Stdin()

	// --- Then ---
	assert.Equal(t, "abc", have)
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
