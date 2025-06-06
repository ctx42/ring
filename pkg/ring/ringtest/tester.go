// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ringtest

import (
	"bytes"
	"maps"

	"github.com/ctx42/testing/pkg/tester"
	"github.com/ctx42/testing/pkg/tstkit"

	"github.com/ctx42/ring/pkg/ring"
)

// Tester represents CLI test helper.
type Tester struct {
	rng  *ring.Ring     // The test ring.
	sin  *bytes.Buffer  // Buffer representing standard input.
	sout *tstkit.Buffer // Buffer to collect stdout writes.
	eout *tstkit.Buffer // Buffer to collect stderr writes.
	t    tester.T       // The test manager.
}

// New returns new instance of [Tester] with given options. By default, the
// constructed [ring.Ring] test instance is returned with:
//
//   - name set to the current program name,
//   - arguments set to empty slice,
//   - environment set to [os.Environ],
//   - metadata set to an empty map,
//   - clock set to [ring.NowUTC],
//   - standard input set to empty [bytes.Buffer],
//   - standard output set to [tstkit.DryBuffer],
//   - standard error set to [tstkit.DryBuffer],
func New(t tester.T, opts ...ring.Option) *Tester {
	t.Helper()
	opts = append([]ring.Option{ring.WithArgs(nil)}, opts...)
	tst := &Tester{
		rng: ring.New(opts...),
		t:   t,
	}
	if tst.sin == nil {
		tst.sin = &bytes.Buffer{}
	}
	if tst.sout == nil {
		tst.sout = tstkit.DryBuffer(t, "stdout")
	}
	if tst.eout == nil {
		tst.eout = tstkit.DryBuffer(t, "stderr")
	}
	return tst
}

// Ring returns a command environment based on [Tester] fields.
func (tst *Tester) Ring(args ...string) *ring.Ring {
	opts := []ring.Option{
		ring.WithEnv(tst.rng.EnvAll()),
		ring.WithMeta(maps.Clone(tst.rng.MetaAll())),
		ring.WithClock(tst.rng.Clock()),
		ring.WithName(tst.rng.Name()),
		ring.WithArgs(args),
	}
	rng := ring.New(opts...)
	rng.SetStdin(tst.sin)
	rng.SetStdout(tst.sout)
	rng.SetStderr(tst.eout)
	return rng
}

// Streams returns standard streams based on [Tester] fields.
func (tst *Tester) Streams() *ring.IO {
	ios := ring.NewIO()
	ios.SetStdin(tst.sin)
	ios.SetStdout(tst.sout)
	ios.SetStderr(tst.eout)
	return ios
}

// SetStdin set buffer to read from as standard input.
func (tst *Tester) SetStdin(sin *bytes.Buffer) *Tester {
	tst.sin = sin
	return tst
}

// WetStdout sets the expectation that standard output will be written to.
func (tst *Tester) WetStdout() *Tester {
	tst.t.Helper()
	tst.sout = tstkit.WetBuffer(tst.t, "stdout")
	return tst
}

// ResetStdout resets the standard output buffer removing all written data.
func (tst *Tester) ResetStdout() { tst.sout.Reset() }

// WetStderr sets expectation that standard error will be written to.
func (tst *Tester) WetStderr() *Tester {
	tst.t.Helper()
	tst.eout = tstkit.WetBuffer(tst.t, "stderr")
	return tst
}

// ResetStderr resets the standard error buffer removing all written data.
func (tst *Tester) ResetStderr() { tst.eout.Reset() }

func (tst *Tester) Stdin() string  { return tst.sin.String() }
func (tst *Tester) Stdout() string { return tst.sout.String() }
func (tst *Tester) Stderr() string { return tst.eout.String() }
