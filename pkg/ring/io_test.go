// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"bytes"
	"os"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NewIO(t *testing.T) {
	// --- When ---
	have := NewIO()

	// --- Then ---
	assert.Same(t, os.Stdin, have.Stdin())
	assert.Same(t, os.Stdout, have.Stdout())
	assert.Same(t, os.Stderr, have.Stderr())
}

func Test_IO_Stdin(t *testing.T) {
	// --- Given ---
	ios := IO{stdin: &bytes.Buffer{}}

	// --- When ---
	have := ios.Stdin()

	// --- Then ---
	assert.Same(t, ios.stdin, have)
}

func Test_IO_Stdout(t *testing.T) {
	// --- Given ---
	ios := IO{stdout: &bytes.Buffer{}}

	// --- When ---
	have := ios.Stdout()

	// --- Then ---
	assert.Same(t, ios.stdout, have)
}

func Test_IO_Stderr(t *testing.T) {
	// --- Given ---
	ios := IO{stderr: &bytes.Buffer{}}

	// --- When ---
	have := ios.Stderr()

	// --- Then ---
	assert.Same(t, ios.stderr, have)
}

func Test_IO_SetStdin(t *testing.T) {
	// --- Given ---
	ios := IO{stdin: &bytes.Buffer{}}

	other := &bytes.Buffer{}

	// --- When ---
	ios.SetStdin(other)

	// --- Then ---
	assert.Same(t, other, ios.stdin)
}

func Test_IO_SetStdout(t *testing.T) {
	// --- Given ---
	ios := IO{stdout: &bytes.Buffer{}}

	other := &bytes.Buffer{}

	// --- When ---
	ios.SetStdout(other)

	// --- Then ---
	assert.Same(t, other, ios.stdout)
}

func Test_IO_SetStderr(t *testing.T) {
	// --- Given ---
	ios := IO{stderr: &bytes.Buffer{}}

	other := &bytes.Buffer{}

	// --- When ---
	ios.SetStderr(other)

	// --- Then ---
	assert.Same(t, other, ios.stderr)
}

func Test_IO_IOClone(t *testing.T) {
	// --- Given ---
	ios := &IO{
		stdin:  &bytes.Buffer{},
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
	}

	// --- When ---
	have := ios.IOClone()

	// --- Then ---
	assert.NotSame(t, ios, have)
	assert.Same(t, ios.stdin, have.stdin)
	assert.Same(t, ios.stdout, have.stdout)
	assert.Same(t, ios.stderr, have.stderr)

	assert.Fields(t, 3, IO{})
}
