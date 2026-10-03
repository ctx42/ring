// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"io"
	"os"
)

// Streamer defines an interface for accessing a program's standard I/O streams.
type Streamer interface {
	// Stdin returns the standard input.
	Stdin() io.Reader

	// Stdout returns the standard output.
	Stdout() io.Writer

	// Stderr returns the standard error.
	Stderr() io.Writer
}

// IO implements Streamer.
var _ Streamer = &IO{}

// IO represents a program's standard I/O streams.
type IO struct {
	stdin  io.Reader // This is the program's standard input.
	stdout io.Writer // This is the program's standard output.
	stderr io.Writer // This is the program's standard error.
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

// SetStdin stores sin, including a nil reader, as the standard input.
func (ios *IO) SetStdin(sin io.Reader) { ios.stdin = sin }

// SetStdout stores sout, including a nil writer, as the standard output.
func (ios *IO) SetStdout(sout io.Writer) { ios.stdout = sout }

// SetStderr stores eout, including a nil writer, as the standard error.
func (ios *IO) SetStderr(eout io.Writer) { ios.stderr = eout }

// IOClone returns a copy of this [IO] with the same streams. The IO prefix
// keeps the promoted method from taking the name [Ring.Clone].
func (ios *IO) IOClone() *IO {
	return &IO{
		stdin:  ios.stdin,
		stdout: ios.stdout,
		stderr: ios.stderr,
	}
}
