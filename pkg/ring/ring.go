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
	// ErrReqMeta indicates a required metadata key is missing.
	ErrReqMeta = errors.New("required ring metadata key")

	// ErrInvMeta indicates a metadata key is invalid due to incorrect type,
	// format, or value.
	ErrInvMeta = errors.New("invalid ring metadata key")
)

// Clock defines a function signature that returns the current time in UTC.
type Clock func() time.Time

// Option configures a [Ring] during creation with [New].
type Option func(*Ring)

// WithEnv configures a [Ring] with the given environment variables.
func WithEnv(env []string) Option {
	return func(rng *Ring) { rng.hidEnv = NewEnv(env) }
}

// WithArgs configures a [Ring] with the given program arguments (excluding
// the program name).
func WithArgs(args []string) Option {
	return func(rng *Ring) { rng.args = args }
}

// WithClock configures a [Ring] with a custom [Clock] function for time.
func WithClock(clk Clock) Option {
	return func(rng *Ring) { rng.clock = clk }
}

// WithMeta configures a [Ring] with the given metadata.
func WithMeta(m map[string]any) Option { return func(rng *Ring) { rng.m = m } }

type hidEnv = Env // Do not export embedded.

var _ Streams = Ring{} // Compile time check.

// Ring represents a program execution context, encapsulating standard I/O
// streams, environment variables, arguments, a clock, and metadata.
type Ring struct {
	hidEnv                // Program environment.
	io     StdIO          // Standard I/O streams.
	clock  Clock          // Function returning current time in UTC.
	args   []string       // Program arguments (excluding program name).
	name   string         // Program name.
	m      map[string]any // Arbitrary metadata.
}

// defaultRing returns a new [Ring] with default configuration.
//
// Configuration:
//   - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
//   - Clock: [Now]
//   - Args: os.Args[1:] (excludes program name)
//   - Name: os.Args[0] (program name)
//   - Environment: nil
//   - Metadata: nil
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

// New creates a new [Ring] with the provided options.
//
// If no options are specified, it defaults to:
//   - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
//   - Environment: [os.Environ]
//   - Clock: [Now]
//   - Args: os.Args[1:]
//   - Name: os.Args[0]
//   - Metadata: empty map
//
// Example:
//
//	rng := New(
//	  WithEnv([]string{"KEY=value"}),
//	  WithArgs([]string{"-a", "arg1"}),
//	)
func New(opts ...Option) Ring {
	rng := defaultRing()
	for _, opt := range opts {
		opt(&rng)
	}
	if rng.hidEnv.e == nil {
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

// Args returns the program arguments, excluding the program name.
func (rng Ring) Args() []string { return rng.args }

// Name returns program name.
func (rng Ring) Name() string { return rng.name }

// WithStdin returns a new [Ring] with the specified standard input stream.
func (rng Ring) WithStdin(sin io.Reader) Ring {
	rng.io.stdin = sin
	return rng
}

// WithStdout returns a new [Ring] with the specified standard output stream.
func (rng Ring) WithStdout(sout io.Writer) Ring {
	rng.io.stdout = sout
	return rng
}

// WithStderr returns a new [Ring] with the specified standard error stream.
func (rng Ring) WithStderr(eout io.Writer) Ring {
	rng.io.stderr = eout
	return rng
}

// WithArgs returns a new [Ring] with the specified program arguments.
func (rng Ring) WithArgs(args []string) Ring {
	rng.args = args
	return rng
}

// WithName returns a new [Ring] with the specified program name.
func (rng Ring) WithName(name string) Ring {
	rng.name = name
	return rng
}

// Env returns a copy of the environment as an [Env].
func (rng Ring) Env() Env { return NewEnv(slices.Clone(rng.EnvGetAll())) }

// EnvSet returns a new [Ring] with the specified environment variable set.
func (rng Ring) EnvSet(key, value string) Ring {
	rng.hidEnv = NewEnv(slices.Clone(rng.EnvGetAll()))
	rng.hidEnv.EnvSet(key, value)
	return rng
}

// EnvSetBulk returns a new [Ring] with the specified environment variables
// appended to the existing environment.
func (rng Ring) EnvSetBulk(env []string) Ring {
	have := append(slices.Clone(rng.EnvGetAll()), env...)
	rng.hidEnv = NewEnv(have)
	return rng
}

// EnvUnset returns a new [Ring] with the specified environment variable
// removed.
func (rng Ring) EnvUnset(key string) Ring {
	rng.hidEnv = NewEnv(slices.Clone(rng.EnvGetAll()))
	rng.hidEnv.EnvUnset(key)
	return rng
}

// MetaSet returns [Ring] with the specified metadata key-value pair set.
// TODO(rz): explain it does not work like env. It is modified in all instances.
func (rng Ring) MetaSet(key string, value any) Ring {
	rng.m[key] = value
	return rng
}

// MetaDelete returns [Ring] with the specified metadata key removed.
// TODO(rz): explain it does not work like env. It is modified in all instances.
func (rng Ring) MetaDelete(key string) Ring {
	delete(rng.m, key)
	return rng
}

// MetaAll returns a clone of the metadata map.
func (rng Ring) MetaAll() map[string]any { return maps.Clone(rng.m) }

// IsZero returns true if the [Ring] is its zero value (e.g., Ring{}), meaning
// it has no environment, metadata, arguments, or name.
func (rng Ring) IsZero() bool {
	if len(rng.hidEnv.e) > 0 {
		return false
	}
	return (rng.m == nil || len(rng.m) == 0) &&
		len(rng.args) == 0 &&
		rng.name == ""
}
