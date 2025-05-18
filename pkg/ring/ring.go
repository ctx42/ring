// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"errors"
	"io"
	"maps"
	"os"
	"slices"
	"time"
)

// Sentinel errors.
var (
	// ErrReqMeta is returned when required [Ring] metadata is messing.
	ErrReqMeta = errors.New("required ring metadata key")

	// ErrInvMeta is returned when [Ring] metadata is invalid either because it
	// is of a wrong type, format or value.
	ErrInvMeta = errors.New("invalid ring metadata key")
)

// Clock is function signature returning current time in UTC.
type Clock func() time.Time

// Option represents [Ring] option.
type Option func(*Ring)

// WithEnv is an option for [New] setting environment.
func WithEnv(env []string) Option {
	return func(rng *Ring) { rng.hidEnv = NewEnv(env) }
}

// WithArgs is an option for [New] setting arguments.
func WithArgs(args []string) Option {
	return func(rng *Ring) { rng.args = args }
}

// WithClock is an option for [New] setting a [Clock] function.
func WithClock(clk Clock) Option {
	return func(rng *Ring) { rng.clock = clk }
}

// WithMeta is an option for [New] setting metadata.
func WithMeta(m map[string]any) Option { return func(rng *Ring) { rng.m = m } }

// Do not export embedded.
type (
	hidEnv = Env
)

var _ Streams = Ring{} // Compile time check.

// Ring represents program execution context.
type Ring struct {
	hidEnv                // Environment.
	io     StdIO          // I/O streams.
	clock  Clock          // Function returning current time in UTC.
	args   []string       // Arguments (without program name).
	name   string         // Context name.
	m      map[string]any // Metadata associated with the ring.
}

// Now returns the current time in UTC.
//
// The only difference between [time.Now] and this function is it always
// returns time in UTC.
func Now() time.Time { return time.Now().UTC() }

// defaultRing returns a new [Ring] with default configuration.
//
// Configuration:
//   - stdin points to [os.Stdin]
//   - stdout points to [os.Stdout]
//   - stderr points to [os.Stderr]
//   - env is nil
//   - meta is nil
//   - args is set as os.Args[1:]
//   - name is set as os.Args[0]
func defaultRing() Ring {
	return Ring{
		io: StdIO{
			stdin:  os.Stdin,
			stdout: os.Stdout,
			stderr: os.Stderr,
		},
		clock: Now,
		name:  os.Args[0],
		args:  os.Args[1:],
	}
}

// New creates new program execution context with the provided options. If no
// options are provided, the configuration is:
//
//   - stdin points to [os.Stdin]
//   - stdout points to [os.Stdout]
//   - stderr points to [os.Stderr]
//   - env is set from [os.Environ]
//   - meta is empty
//   - args set as os.Args[1:]
//   - name set as os.Args[0]
func New(opts ...Option) Ring {
	rng := defaultRing()
	for _, opt := range opts {
		opt(&rng)
	}
	if rng.EnvIsNil() {
		rng.hidEnv = NewEnv(os.Environ())
	}
	if rng.m == nil {
		rng.m = make(map[string]any, 10)
	}
	return rng
}

// Stdin returns the standard input.
func (rng Ring) Stdin() io.Reader { return rng.io.stdin }

// Stdout returns the standard output.
func (rng Ring) Stdout() io.Writer { return rng.io.stdout }

// Stderr returns the standard error.
func (rng Ring) Stderr() io.Writer { return rng.io.stderr }

// Clock returns function returning current time in UTC.
func (rng Ring) Clock() func() time.Time { return rng.clock }

// Args returns the arguments program was started with, without a program name.
func (rng Ring) Args() []string { return rng.args }

// Name returns program name.
func (rng Ring) Name() string { return rng.name }

// WithStdin returns new [Ring] with given standard input.
func (rng Ring) WithStdin(sin io.Reader) Ring {
	rng.io.stdin = sin
	return rng
}

// WithStdout returns new [Ring] with given standard output.
func (rng Ring) WithStdout(sout io.Writer) Ring {
	rng.io.stdout = sout
	return rng
}

// WithStderr returns new [Ring] with given standard error.
func (rng Ring) WithStderr(eout io.Writer) Ring {
	rng.io.stderr = eout
	return rng
}

// WithArgs returns new [Ring] with given program arguments.
func (rng Ring) WithArgs(args []string) Ring {
	rng.args = args
	return rng
}

// WithName returns new [Ring] with program name.
func (rng Ring) WithName(name string) Ring {
	rng.name = name
	return rng
}

// Env returns copy of the environment.
func (rng Ring) Env() Env { return NewEnv(slices.Clone(rng.EnvGetAll())) }

// EnvSet returns new [Ring] with a given environment variable set.
func (rng Ring) EnvSet(key, value string) Ring {
	have := slices.Clone(rng.EnvGetAll())
	rng.hidEnv = NewEnv(have)
	rng.hidEnv.EnvSet(key, value)
	return rng
}

// EnvSetBulk returns new [Ring] with environment added.
func (rng Ring) EnvSetBulk(env []string) Ring {
	have := slices.Clone(rng.EnvGetAll())
	have = append(have, env...)
	rng.hidEnv = NewEnv(have)
	return rng
}

// EnvUnset returns new [Ring] with given environment variable deleted.
func (rng Ring) EnvUnset(key string) Ring {
	have := slices.Clone(rng.EnvGetAll())
	rng.hidEnv = NewEnv(have)
	rng.hidEnv.EnvUnset(key)
	return rng
}

// MetaSet returns new [Ring] with metadata variable named by key set.
func (rng Ring) MetaSet(key string, value any) Ring {
	rng.m[key] = value
	return rng
}

// MetaDelete returns new [Ring] with metadata variable named by key deleted.
func (rng Ring) MetaDelete(key string) Ring {
	delete(rng.m, key)
	return rng
}

// MetaAll returns clone of metadata associated with the ring.
func (rng Ring) MetaAll() map[string]any { return maps.Clone(rng.m) }

// IsZero returns true when ring is zero value - `ring.Ring{}` for example.
func (rng Ring) IsZero() bool {
	return rng.EnvIsNil() &&
		(rng.m == nil || len(rng.m) == 0) &&
		len(rng.args) == 0 &&
		rng.name == ""
}
