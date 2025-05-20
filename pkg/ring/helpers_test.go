// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_NowUTC(t *testing.T) {
	// --- When ---
	have := NowUTC()

	// --- Then ---
	assert.Within(t, time.Now(), "1ms", have)
	assert.Zone(t, time.UTC, have.Location())
}
