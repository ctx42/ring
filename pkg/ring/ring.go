// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

// Package ring provides a program execution context that bundles standard
// I/O, environment variables, arguments, a clock, a filesystem, and
// metadata.
//
// Import it as "github.com/ctx42/ring/pkg/ring".
package ring

import (
	"errors"
	"io/fs"
	"os"
	"slices"
	"time"
)

// Sentinel errors.
var (
	// ErrReqMeta is the sentinel callers return when a required metadata
	// key is missing.
	ErrReqMeta = errors.New("required ring metadata key")

	// ErrInvMeta is the sentinel callers return when a metadata value has
	// the wrong type, format, or value.
	ErrInvMeta = errors.New("invalid ring metadata key")

	// ErrNoFsAccess is returned when [Ring] has no filesystem access.
	ErrNoFsAccess = errors.New("no filesystem access")
)

// Clock defines a function signature that returns the current time in UTC.
type Clock func() time.Time

// Option configures a [Ring] during creation with [New].
type Option func(*Ring)

// WithEnv configures a [Ring] with the given environment variables. A nil
// slice sets an empty environment. Omit WithEnv to keep [os.Environ], which
// [New] applies when no environment was configured.
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

// WithFS configures a [Ring] with access to a read-only filesystem.
func WithFS(filesystem fs.FS) Option {
	return func(rng *Ring) { rng.fs = filesystem }
}

// Unexported aliases hide the embedded field names.
type (
	hidEnv = Env
	hidIO  = IO
)

// Ring implements Streamer.
var _ Streamer = Ring{}

// Ring implements Environ.
var _ Environ = Ring{}

// Ring represents a program execution context, encapsulating standard I/O
// streams, environment variables, arguments, a clock, a filesystem, and
// metadata.
//
// Do not use the zero value. Its promoted environment and stream methods
// panic, and [Ring.Clock] is nil. Build rings with [New]. [Ring.Clone] of a
// zero ring keeps the nil environment and the nil streams.
type Ring struct {
	*hidEnv                // This is the program environment.
	*hidIO                 // These are the standard I/O streams.
	clock   Clock          // This returns the current time in UTC.
	fs      fs.FS          // This is the program filesystem.
	name    string         // This is the program name.
	args    []string       // These are the arguments, excluding the name.
	meta    map[string]any // This is arbitrary metadata.
}

// defaultRing returns a new [Ring] with default configuration.
//
// Configuration:
//   - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
//   - Clock: [NowUTC]
//   - Name: os.Args[0] when present, otherwise empty
//   - Args: os.Args[1:] when os.Args is non-empty, otherwise nil
//   - Environment: nil
//   - Metadata: nil
//   - Filesystem: nil
func defaultRing() *Ring {
	var name string
	var args []string
	if len(os.Args) > 0 {
		name = os.Args[0]
		args = os.Args[1:]
	}
	return &Ring{
		hidIO: NewIO(),
		clock: NowUTC,
		name:  name,
		args:  args,
	}
}

// New creates a new [Ring] with the provided options.
//
// If no options are specified, it defaults to:
//   - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
//   - Environment: [os.Environ]
//   - Clock: [NowUTC], including when a clock option is nil
//   - Args: os.Args[1:] when os.Args is non-empty, otherwise nil
//   - Name: os.Args[0] when present, otherwise empty
//   - Metadata: empty map
//   - Filesystem: no access.
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
		if opt == nil {
			continue
		}
		opt(rng)
	}
	if rng.hidEnv == nil {
		rng.hidEnv = NewEnv(os.Environ())
	}
	if rng.meta == nil {
		rng.meta = make(map[string]any)
	}
	if rng.clock == nil {
		rng.clock = NowUTC
	}
	return rng
}

// Clock returns the function returning the current time in UTC.
func (rng *Ring) Clock() func() time.Time { return rng.clock }

// Args returns the program arguments, excluding the program name. The
// returned slice is not a copy; callers must not modify it.
func (rng *Ring) Args() []string { return rng.args }

// SetArgs sets the program arguments, excluding the program name.
func (rng *Ring) SetArgs(args []string) *Ring {
	rng.args = args
	return rng
}

// Name returns the program name.
func (rng *Ring) Name() string { return rng.name }

// MetaSet sets the metadata value for the given key. If the key already exists,
// its value is overwritten. The value may be any type, including nil.
func (rng *Ring) MetaSet(key string, value any) {
	if rng.meta == nil {
		rng.meta = make(map[string]any)
	}
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

// MetaAll returns the live metadata map. The map is shared across all
// [Ring.Clone] copies; mutations are visible to every clone.
func (rng *Ring) MetaAll() map[string]any { return rng.meta }

// FS returns a hierarchical file system associated with the instance.
func (rng *Ring) FS() (fs.FS, error) {
	if rng.fs == nil {
		return nil, ErrNoFsAccess
	}
	return rng.fs, nil
}

// Clone copies the [Ring]. The copy has its own environment map and its
// own argument slice. Metadata, the filesystem, and the standard streams
// are shared, so a metadata change or a write to a stream is visible on
// every clone.
func (rng *Ring) Clone() *Ring {
	var env *Env
	if rng.hidEnv != nil {
		env = rng.hidEnv.EnvClone()
	}
	var ios *IO
	if rng.hidIO != nil {
		ios = rng.hidIO.IOClone()
	}
	return &Ring{
		hidEnv: env,
		hidIO:  ios,
		clock:  rng.clock,
		fs:     rng.fs,
		name:   rng.name,
		args:   slices.Clone(rng.args),
		meta:   rng.meta,
	}
}
