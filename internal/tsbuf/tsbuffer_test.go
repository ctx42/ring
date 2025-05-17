// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package tsbuf

import (
	"testing"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
	"github.com/ctx42/testing/pkg/tester"
)

func Test_NewTSBuffer(t *testing.T) {
	t.Run("without name", func(t *testing.T) {
		// --- When ---
		tsb := NewTSBuffer()

		// --- Then ---
		assert.Empty(t, tsb.name)
		assert.Equal(t, TSBuffDefault, tsb.kind)
		assert.Empty(t, tsb.buf.String())
		assert.True(t, tsb.check)
		assert.Equal(t, 0, tsb.wc)
		assert.Equal(t, 0, tsb.rc)
	})

	t.Run("with name", func(t *testing.T) {
		// --- When ---
		tsb := NewTSBuffer("name")

		// --- Then ---
		assert.Equal(t, "name ", tsb.name)
		assert.Empty(t, tsb.buf.String())
		assert.True(t, tsb.check)
		assert.Equal(t, 0, tsb.wc)
		assert.Equal(t, 0, tsb.rc)
	})
}

func Test_TSBuffer_Name(t *testing.T) {
	t.Run("without name", func(t *testing.T) {
		// --- When ---
		tsb := NewTSBuffer()

		// --- Then ---
		assert.Empty(t, tsb.Name())
	})

	t.Run("with name", func(t *testing.T) {
		// --- When ---
		tsb := NewTSBuffer("name")

		// --- Then ---
		assert.Equal(t, "name", tsb.Name())
	})
}

func Test_TSBuffer_Type(t *testing.T) {
	// --- Given ---
	tsb := NewTSBuffer()

	// --- When ---
	have := tsb.Kind()

	// --- Then ---
	assert.Equal(t, TSBuffDefault, have)
}

func Test_TSBuffer_SkipChecks(t *testing.T) {
	// --- Given ---
	tsb := NewTSBuffer()

	// --- When ---
	tsb.SkipChecks()

	// --- Then ---
	assert.False(t, tsb.check)
}

func Test_TSBuffer_Write_String(t *testing.T) {
	// --- Given ---
	tsb := NewTSBuffer()

	// --- When ---
	n, err := tsb.Write([]byte{97, 98, 99})

	// --- Then ---
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 1, tsb.wc)
	assert.Equal(t, 0, tsb.rc)
	assert.Equal(t, "abc", tsb.String())
	assert.Equal(t, 1, tsb.wc)
	assert.Equal(t, 1, tsb.rc)
}

func Test_TSBuffer_Write_string(t *testing.T) {
	t.Run("do not increase read counter", func(t *testing.T) {
		// --- Given ---
		tsb := NewTSBuffer()
		must.Value(tsb.Write([]byte{97, 98, 99}))
		_ = tsb.String()

		// --- When ---
		have := tsb.string(false)

		// --- Then ---
		assert.Equal(t, 1, tsb.wc)
		assert.Equal(t, 1, tsb.rc)
		assert.Equal(t, "abc", have)
	})

	t.Run("increase read counter", func(t *testing.T) {
		// --- Given ---
		tsb := NewTSBuffer()
		must.Value(tsb.Write([]byte{97, 98, 99}))
		_ = tsb.String()

		// --- When ---
		have := tsb.string(true)

		// --- Then ---
		assert.Equal(t, 1, tsb.wc)
		assert.Equal(t, 2, tsb.rc)
		assert.Equal(t, "abc", have)
	})
}

func Test_TSBuffer_WriteString_String(t *testing.T) {
	// --- Given ---
	tsb := NewTSBuffer()

	// --- When ---
	n, err := tsb.WriteString("abc")

	// --- Then ---
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, 1, tsb.wc)
	assert.Equal(t, 0, tsb.rc)
	assert.Equal(t, "abc", tsb.String())
	assert.Equal(t, 1, tsb.wc)
	assert.Equal(t, 1, tsb.rc)
}

func Test_TSBuffer_MustWriteString(t *testing.T) {
	// --- Given ---
	tsb := NewTSBuffer()

	// --- When ---
	tsb.MustWriteString("abc")

	// --- Then ---
	assert.Equal(t, 1, tsb.wc)
	assert.Equal(t, 0, tsb.rc)
	assert.Equal(t, "abc", tsb.String())
}

func Test_TSBuffer_Reset(t *testing.T) {
	// --- Given ---
	tsb := NewTSBuffer()
	must.Value(tsb.WriteString("abc"))
	_ = tsb.String()

	// --- When ---
	tsb.Reset()

	// --- Then ---
	assert.Equal(t, 0, tsb.wc)
	assert.Equal(t, 0, tsb.rc)
	assert.Equal(t, "", tsb.string(false))
}

func Test_DryBuffer(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)

		// --- Then ---
		assert.Equal(t, TSBuffDry, buf.Kind())
	})

	t.Run("buf not written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)

		// --- Then ---
		assert.NotNil(t, buf)
	})

	t.Run("buf written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer to be empty:\n" +
			"\twant: \n" +
			"\thave: \"abc\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("buf written checks skipped", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy)
		buf.SkipChecks()
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("buf named and written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buf-name buffer to be empty:\n" +
			"\twant: \n" +
			"\thave: \"abc\""
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := DryBuffer(tspy, "buf-name ")
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})
}

func Test_WetBuffer(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)

		// --- Then ---
		assert.Equal(t, TSBuffWet, buf.Kind())
		buf.SkipChecks()
	})

	t.Run("buf written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "abc", buf.String())
	})

	t.Run("buf written but not examined", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer to be examined with String or Buffer methods"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("buf written but not examined checks skipped", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)
		buf.SkipChecks()
		_, err := buf.WriteString("abc")

		// --- Then ---
		assert.NoError(t, err)
		assert.NotNil(t, buf)
	})

	t.Run("buf not written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buffer not to be empty"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy)

		// --- Then ---
		assert.NotNil(t, buf)
	})

	t.Run("buf named and not written", func(t *testing.T) {
		// --- Given ---
		tspy := tester.New(t, 2)
		tspy.ExpectCleanups(1)
		tspy.ExpectError()
		wMsg := "expected buf-name buffer not to be empty"
		tspy.ExpectLogEqual(wMsg)
		tspy.Close()

		// --- When ---
		buf := WetBuffer(tspy, "buf-name ")

		// --- Then ---
		assert.NotNil(t, buf)
	})
}
