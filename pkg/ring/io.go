// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"io"
	"os"
)

// Streamer defines an interface for accessing a program's standard I/O streams.
type Streamer interface {
	Stdin() io.Reader  // Standard input.
	Stdout() io.Writer // Standard output.
	Stderr() io.Writer // Standard error.
}

var _ Streamer = &IO{} // Compile time check.

// IO represents a program's standard I/O streams.
type IO struct {
	stdin  io.Reader // Program standard input.
	stdout io.Writer // Program standard output.
	stderr io.Writer // Program standard error.
}

// NewIO returns a new instance of [IO] with [os.Stdin], [os.Stdout],
// and [os.Stderr] as default values for the stdin, stdout, and stderr fields
// respectively.
func NewIO() *IO {
	return &IO{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (ios *IO) Stdin() io.Reader { return ios.stdin }

func (ios *IO) Stdout() io.Writer { return ios.stdout }

func (ios *IO) Stderr() io.Writer { return ios.stderr }

// SetStdin sets the standard input stream.
func (ios *IO) SetStdin(sin io.Reader) { ios.stdin = sin }

// SetStdout sets the standard output stream.
func (ios *IO) SetStdout(sout io.Writer) { ios.stdout = sout }

// SetStderr sets the standard error stream.
func (ios *IO) SetStderr(eout io.Writer) { ios.stderr = eout }

// IOClone creates a copy of the current [IO] instance with identical streams.
func (ios *IO) IOClone() *IO {
	return &IO{
		stdin:  ios.stdin,
		stdout: ios.stdout,
		stderr: ios.stderr,
	}
}
