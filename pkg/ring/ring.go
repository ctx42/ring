// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"errors"
	"os"
	"time"
)

// Sentinel errors.
var (
	// ErrReqMeta indicates a required metadata key is missing.
	ErrReqMeta = errors.New("required ring metadata key")

	// ErrInvMeta indicates a metadata key is invalid due to an incorrect type,
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

// WithName configures a [Ring] with the given program name.
func WithName(name string) Option {
	return func(rng *Ring) { rng.name = name }
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
func WithMeta(meta map[string]any) Option {
	return func(rng *Ring) { rng.meta = meta }
}

// Hide embedded fields.
type (
	hidEnv = Env
	hidIO  = IO
)

var _ Streamer = Ring{} // Compile time check.

// Ring represents a program execution context, encapsulating standard I/O
// streams, environment variables, arguments, a clock, and metadata.
type Ring struct {
	*hidEnv                // Program environment.
	*hidIO                 // Standard I/O streams.
	clock   Clock          // Function returning current time in UTC.
	name    string         // Program name.
	args    []string       // Program arguments (excluding program name).
	meta    map[string]any // Arbitrary metadata.
}

// defaultRing returns a new [Ring] with default configuration.
//
// Configuration:
//   - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
//   - Clock: [NowUTC]
//   - Name: os.Args[0] (program name)
//   - Args: os.Args[1:] (excludes program name)
//   - Environment: nil
//   - Metadata: nil
func defaultRing() *Ring {
	return &Ring{
		hidIO: NewIO(),
		clock: NowUTC,
		name:  os.Args[0],
		args:  os.Args[1:],
	}
}

// New creates a new [Ring] with the provided options.
//
// If no options are specified, it defaults to:
//   - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
//   - Environment: [os.Environ]
//   - Clock: [NowUTC]
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
func New(opts ...Option) *Ring {
	rng := defaultRing()
	for _, opt := range opts {
		opt(rng)
	}
	if rng.hidEnv == nil {
		rng.hidEnv = NewEnv(os.Environ())
	}
	if rng.meta == nil {
		rng.meta = make(map[string]any)
	}
	return rng
}

// Clock returns function returning current time in UTC.
func (rng *Ring) Clock() func() time.Time { return rng.clock }

// Args returns the program arguments, excluding the program name.
func (rng *Ring) Args() []string { return rng.args }

// SetArgs sets the program arguments, excluding the program name.
func (rng *Ring) SetArgs(args []string) *Ring {
	rng.args = args
	return rng
}

// Name returns program name.
func (rng *Ring) Name() string { return rng.name }

// MetaSet sets the metadata value for the given key. If the key already exists,
// its value is overwritten. The value may be any type, including nil.
func (rng *Ring) MetaSet(key string, value any) {
	rng.meta[key] = value
}

// MetaGet retrieves the metadata value associated with the given key. If the
// key exists, it returns the value, which may be nil or empty. If the key does
// not exist, it returns nil.
func (rng *Ring) MetaGet(key string) any {
	return rng.meta[key]
}

// MetaLookup retrieves the metadata value associated with the given key. If
// the key exists in the metadata, it returns the value (which may be nil or
// empty) and true. If the key does not exist, it returns nil and false.
func (rng *Ring) MetaLookup(key string) (any, bool) {
	val, ok := rng.meta[key]
	return val, ok
}

// MetaDelete removes the metadata value associated with the given key. If the
// key does not exist, the method has no effect.
func (rng *Ring) MetaDelete(key string) {
	delete(rng.meta, key)
}

// MetaAll returns metadata map.
func (rng *Ring) MetaAll() map[string]any { return rng.meta }
