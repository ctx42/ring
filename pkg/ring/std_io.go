// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"io"
	"os"
)

// Streams represents standard I/O streams.
type Streams interface {
	Stdin() io.Reader  // Standard input.
	Stdout() io.Writer // Standard output.
	Stderr() io.Writer // Standard error.
}

var _ Streams = StdIO{} // Compile time check.

// StdIO represents program I/O streams.
type StdIO struct {
	stdin  io.Reader // Program standard input.
	stdout io.Writer // Program standard output.
	stderr io.Writer // Program standard error.
}

// NewStdIO returns a new instance of the StdIO struct with os.Stdin,
// os.Stdout, and os.Stderr as default values for the stdin, stdout, and stderr
// fields respectively.
func NewStdIO() Streams {
	return StdIO{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// Stdin returns the standard input to use for a program.
func (sio StdIO) Stdin() io.Reader { return sio.stdin }

// Stdout returns the standard output to use for a program.
func (sio StdIO) Stdout() io.Writer { return sio.stdout }

// Stderr returns the standard error to use for a program.
func (sio StdIO) Stderr() io.Writer { return sio.stderr }

// WithStdin returns [StdIO] with given standard input.
func (sio StdIO) WithStdin(sin io.Reader) StdIO {
	sio.stdin = sin
	return sio
}

// WithStdout returns [StdIO] with given standard output.
func (sio StdIO) WithStdout(sout io.Writer) StdIO {
	sio.stdout = sout
	return sio
}

// WithStderr returns [StdIO] with given standard error.
func (sio StdIO) WithStderr(eout io.Writer) StdIO {
	sio.stderr = eout
	return sio
}
