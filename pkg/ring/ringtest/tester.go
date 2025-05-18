// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ringtest

import (
	"bytes"
	"io"
	"maps"
	"os"
	"slices"
	"sort"
	"time"

	"github.com/ctx42/testing/pkg/tester"
	"github.com/ctx42/testing/pkg/tstkit"

	"github.com/ctx42/ring/pkg/ring"
)

// WithEnv is an option for [New] setting the environment to use.
func WithEnv(env []string) func(*Tester) {
	return func(tst *Tester) { tst.hidEnv = ring.NewEnv(env) }
}

// WithName is an option for [New] setting program name.
func WithName(name string) func(*Tester) {
	return func(tst *Tester) { tst.name = name }
}

// WithMeta is an option for [New] setting ring metadata.
func WithMeta(m map[string]any) func(*Tester) {
	return func(tst *Tester) { tst.m = m }
}

// WithClock is an option for [New] setting clock function.
func WithClock(clock ring.Clock) func(*Tester) {
	return func(tst *Tester) { tst.clock = clock }
}

// Sort works like [sort.Strings] but returns sorted slice.
func Sort(in []string) []string {
	sort.Strings(in)
	return in
}

// Hide embedded fields.
type (
	hidEnv = ring.Env
)

// Tester represents CLI test helper.
type Tester struct {
	hidEnv                // Environment.
	m      map[string]any // Metadata.
	sin    *bytes.Buffer  // Buffer representing standard input.
	sout   *tstkit.Buffer // Buffer to collect stdout writes.
	eout   *tstkit.Buffer // Buffer to collect stderr writes.
	clock  ring.Clock     // Function returning current time in UTC.
	name   string         // Program name.
	t      tester.T       // The test manager.
}

// New returns new instance of Tester.
func New(t tester.T, opts ...func(*Tester)) *Tester {
	t.Helper()
	tst := &Tester{
		sin:   &bytes.Buffer{},
		sout:  tstkit.DryBuffer(t, "stdout"),
		eout:  tstkit.DryBuffer(t, "stderr"),
		clock: ring.Now,
		name:  os.Args[0],
		t:     t,
	}
	for _, opt := range opts {
		opt(tst)
	}
	if tst.EnvIsNil() {
		tst.hidEnv = ring.NewEnv(os.Environ())
	}
	if tst.m == nil {
		tst.m = make(map[string]any, 10)
	}
	return tst
}

// Ring returns a command environment based on [Tester] fields.
func (tst *Tester) Ring(args ...string) ring.Ring {
	opts := []ring.Option{
		ring.WithEnv(slices.Clone(tst.EnvGetAll())),
		ring.WithMeta(maps.Clone(tst.m)),
		ring.WithClock(tst.clock),
		ring.WithArgs(args),
	}
	return ring.New(opts...).
		WithName(tst.name).
		WithStdin(tst.sin).
		WithStdout(tst.sout).
		WithStderr(tst.eout)
}

// Streams returns standard streams based on [Tester] fields.
func (tst *Tester) Streams() ring.Streams {
	return ring.StdIO{}.
		WithStdin(tst.sin).
		WithStdout(tst.sout).
		WithStderr(tst.eout)
}

// WithStdin set buffer to read from as standard input.
func (tst *Tester) WithStdin(sin *bytes.Buffer) *Tester {
	tst.sin = sin
	return tst
}

// WetStdout sets expectation that standard output will be written to.
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

func (tst *Tester) Name() string            { return tst.name }
func (tst *Tester) Stdin() io.Reader        { return tst.sin }
func (tst *Tester) Stdout() string          { return tst.sout.String() }
func (tst *Tester) Stderr() string          { return tst.eout.String() }
func (tst *Tester) Clock() func() time.Time { return tst.clock }
