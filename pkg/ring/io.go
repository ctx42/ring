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

// IO represents program standard I/O streams.
type IO struct {
	stdin  io.Reader // Program standard input.
	stdout io.Writer // Program standard output.
	stderr io.Writer // Program standard error.
}

// NewIO returns a new instance of the IO struct with [os.Stdin], [os.Stdout],
// and [os.Stderr] as default values for the stdin, stdout, and stderr fields
// respectively.
func NewIO() *IO {
	return &IO{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// Stdin returns the standard input to use for a program.
func (ios *IO) Stdin() io.Reader { return ios.stdin }

// Stdout returns the standard output to use for a program.
func (ios *IO) Stdout() io.Writer { return ios.stdout }

// Stderr returns the standard error to use for a program.
func (ios *IO) Stderr() io.Writer { return ios.stderr }

// SetStdin returns [IO] with the given standard input.
func (ios *IO) SetStdin(sin io.Reader) { ios.stdin = sin }

// SetStdout returns [IO] with the given standard output.
func (ios *IO) SetStdout(sout io.Writer) { ios.stdout = sout }

// SetStderr returns [IO] with the given standard error.
func (ios *IO) SetStderr(eout io.Writer) { ios.stderr = eout }
