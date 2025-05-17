// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"bytes"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewStdIO(t *testing.T) {
	// --- When ---
	sio := NewStdIO()

	// --- Then ---
	assert.Same(t, os.Stdin, sio.Stdin())
	assert.Same(t, os.Stdout, sio.Stdout())
	assert.Same(t, os.Stderr, sio.Stderr())
}

func Test_StdIO(t *testing.T) {
	// --- Given ---
	sio := StdIO{
		stdin:  &bytes.Buffer{},
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
	}

	// --- Then ---
	assert.Same(t, sio.stdin, sio.Stdin())
	assert.Same(t, sio.stdout, sio.Stdout())
	assert.Same(t, sio.stderr, sio.Stderr())
}

func Test_StdIO_WithStdin_Stdin(t *testing.T) {
	// --- Given ---
	sio := StdIO{
		stdin: &bytes.Buffer{},
	}
	other := &bytes.Buffer{}

	// --- When ---
	have := sio.WithStdin(other)

	// --- Then ---
	assert.Same(t, other, have.Stdin())
}

func Test_StdIO_WithStdout_Stdout(t *testing.T) {
	// --- Given ---
	sio := StdIO{
		stdout: &bytes.Buffer{},
	}
	other := &bytes.Buffer{}

	// --- When ---
	have := sio.WithStdout(other)

	// --- Then ---
	assert.Same(t, other, have.Stdout())
}

func Test_StdIO_WithStderr_Stderr(t *testing.T) {
	// --- Given ---
	sio := StdIO{
		stderr: &bytes.Buffer{},
	}
	other := &bytes.Buffer{}

	// --- When ---
	have := sio.WithStderr(other)

	// --- Then ---
	assert.Same(t, other, have.Stderr())
}
