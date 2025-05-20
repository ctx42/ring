// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ringtest

import (
	"bytes"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
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
		tst := New(tspy)

		// --- Then ---

		// The instance of [ring.Ring].
		assert.Equal(t, Sort(os.Environ()), Sort(tst.rng.EnvAll()))
		assert.NotNil(t, tst.rng.MetaAll())
		assert.Empty(t, tst.rng.MetaAll())
		assert.Same(t, os.Stdin, tst.rng.Stdin())
		assert.Same(t, os.Stdout, tst.rng.Stdout())
		assert.Same(t, os.Stderr, tst.rng.Stderr())
		assert.Same(t, ring.NowUTC, tst.rng.Clock())
		assert.Equal(t, os.Args[0], tst.rng.Name())
		assert.Empty(t, tst.rng.Args())

		// The instance of [tester.Tester].
		assert.Empty(t, tst.sin.String())
		assert.Equal(t, "", tst.sout.String())
		assert.Equal(t, "", tst.eout.String())
		assert.Same(t, tspy, tst.t)
	})

	t.Run("with environment", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		env := []string{"A=B", "C=D"}

		// --- When ---
		tst := New(tspy, ring.WithEnv(env))

		// --- Then ---
		assert.Equal(t, env, Sort(tst.Ring().EnvAll()))
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
		rng := tst.Ring()

		// --- Then ---
		assert.Equal(t, Sort(os.Environ()), Sort(rng.EnvAll()))
		assert.NotNil(t, rng.MetaAll())
		assert.Empty(t, rng.MetaAll())
		assert.Same(t, tst.sin, rng.Stdin())
		assert.Same(t, tst.sout, rng.Stdout())
		assert.Same(t, tst.eout, rng.Stderr())
		assert.Same(t, ring.NowUTC, rng.Clock())
		assert.Equal(t, os.Args[0], rng.Name())
		assert.Empty(t, rng.Args())
	})

	t.Run("with args", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy)

		// --- When ---
		rng := tst.Ring("a", "b", "c")

		// --- Then ---
		assert.Equal(t, Sort(os.Environ()), Sort(rng.EnvAll()))
		assert.NotNil(t, rng.MetaAll())
		assert.Empty(t, rng.MetaAll())
		assert.Same(t, tst.sin, rng.Stdin())
		assert.Same(t, tst.sout, rng.Stdout())
		assert.Same(t, tst.eout, rng.Stderr())
		assert.Same(t, ring.NowUTC, rng.Clock())
		assert.Equal(t, os.Args[0], rng.Name())
		assert.Equal(t, []string{"a", "b", "c"}, rng.Args())
	})

	t.Run("with a custom name", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t)
		tspy.ExpectCleanups(2)
		tspy.Close()

		tst := New(tspy, ring.WithName("my"))

		// --- When ---
		rng := tst.Ring()

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
	_, _ = tst.sout.WriteString("test")

	// --- When ---
	tst.ResetStdout()

	// --- Then ---
	assert.Equal(t, "", tst.sout.String())
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
