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
	ios := NewIO()

	// --- Then ---
	assert.Same(t, os.Stdin, ios.Stdin())
	assert.Same(t, os.Stdout, ios.Stdout())
	assert.Same(t, os.Stderr, ios.Stderr())
}

func Test_IO_Stdin(t *testing.T) {
	// --- Given ---
	ios := IO{stdin: &bytes.Buffer{}}

	// --- Then ---
	assert.Same(t, ios.stdin, ios.Stdin())
}

func Test_IO_Stdout(t *testing.T) {
	// --- Given ---
	ios := IO{stdout: &bytes.Buffer{}}

	// --- Then ---
	assert.Same(t, ios.stdout, ios.Stdout())
}

func Test_IO_Stderr(t *testing.T) {
	// --- Given ---
	ios := IO{stderr: &bytes.Buffer{}}

	// --- Then ---
	assert.Same(t, ios.stderr, ios.Stderr())
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
