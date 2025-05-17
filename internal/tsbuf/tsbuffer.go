// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package tsbuf

import (
	"bytes"
	"strings"
	"sync"

	"github.com/ctx42/testing/pkg/tester"
)

// TSBuffer kinds.
const (
	TSBuffDry     = "dry"     // Must never be written to.
	TSBuffWet     = "wet"     // Must be written to and check the contents.
	TSBuffDefault = "default" // Created by [NewTSBuffer].
)

// TSBuffer represents thread safe io.Writer.
type TSBuffer struct {
	name  string        // Buffer name.
	kind  string        // Buffer kind (default: TSBuffDefault).
	buf   *bytes.Buffer // Buffer to write to.
	mx    sync.Mutex    // Guards the buffer.
	check bool          // Run cleanups (default: true).
	wc    int           // Write count.
	rc    int           // Read count.
}

// NewTSBuffer returns new instance of TSBuffer. You may provide a name for
// the buffer.
func NewTSBuffer(names ...string) *TSBuffer {
	tsb := &TSBuffer{
		kind:  TSBuffDefault,
		buf:   &bytes.Buffer{},
		mx:    sync.Mutex{},
		check: true,
	}
	if len(names) > 0 {
		tsb.name = strings.TrimSpace(names[0]) + " "
	}
	return tsb
}

// Name returns name of the buffer or empty sting if name was not provided.
func (tsb *TSBuffer) Name() string {
	return strings.TrimSpace(tsb.name)
}

// Kind returns buffer kind. Kind describes how the buffer behaves. Example b
// buffer kinds are buffers created with:
//
//   - [NewTSBuffer]
//   - [DryBuffer]
//   - [WetBuffer]
func (tsb *TSBuffer) Kind() string {
	return tsb.kind
}

// SkipChecks skip test case cleanup checks.
func (tsb *TSBuffer) SkipChecks() {
	tsb.mx.Lock()
	defer tsb.mx.Unlock()
	tsb.check = false
}

func (tsb *TSBuffer) Write(p []byte) (n int, err error) {
	tsb.mx.Lock()
	defer tsb.mx.Unlock()
	tsb.wc++
	return tsb.buf.Write(p)
}

func (tsb *TSBuffer) WriteString(s string) (n int, err error) {
	tsb.mx.Lock()
	defer tsb.mx.Unlock()
	tsb.wc++
	return tsb.buf.WriteString(s)
}

// MustWriteString writes string to the buffer. Panics on error.
func (tsb *TSBuffer) MustWriteString(s string) {
	if _, err := tsb.WriteString(s); err != nil {
		panic(err)
	}
}

func (tsb *TSBuffer) String() string {
	tsb.mx.Lock()
	defer tsb.mx.Unlock()
	return tsb.string(true)
}

// string returns data written to the buffer. Updates read counter if inc is
// true. Assumes the locks were acquired by the caller.
func (tsb *TSBuffer) string(inc bool) string {
	if inc {
		tsb.rc++
	}
	return tsb.buf.String()
}

// Reset resets underlying bytes.Buffer and counters.
func (tsb *TSBuffer) Reset() {
	tsb.mx.Lock()
	defer tsb.mx.Unlock()
	tsb.wc = 0
	tsb.rc = 0
	tsb.buf.Reset()
}

// DryBuffer returns thread save buffer which checks nothing has been written
// to it at the test end. On error, it marks the test as failed.
func DryBuffer(t tester.T, names ...string) *TSBuffer {
	t.Helper()
	tsb := NewTSBuffer(names...)
	tsb.kind = TSBuffDry
	t.Cleanup(func() {
		t.Helper()
		tsb.mx.Lock()
		defer tsb.mx.Unlock()
		if !tsb.check {
			return
		}
		out := tsb.string(false)
		if out != "" {
			format := "expected %sbuffer to be empty:\n" +
				"\twant: \n" +
				"\thave: %q"
			t.Errorf(format, tsb.name, out)
		}
	})
	return tsb
}

// WetBuffer returns thread save buffer which checks something has been written
// to it at the test end. Also checks the content of the buffer was examined by
// calling String method. On error, it marks the test as failed.
func WetBuffer(t tester.T, names ...string) *TSBuffer {
	t.Helper()
	tsb := NewTSBuffer(names...)
	tsb.kind = TSBuffWet
	t.Cleanup(func() {
		t.Helper()
		tsb.mx.Lock()
		defer tsb.mx.Unlock()
		if !tsb.check {
			return
		}
		out := tsb.string(false)
		if out == "" {
			format := "expected %sbuffer not to be empty"
			t.Errorf(format, tsb.name)
			return
		}
		if tsb.rc == 0 {
			format := "expected %sbuffer to be examined " +
				"with String or Buffer methods"
			t.Errorf(format, tsb.name)
		}
	})
	return tsb
}
