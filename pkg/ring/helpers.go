// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"time"
)

// Now returns the current time in UTC, equivalent to [time.Now].UTC().
func Now() time.Time { return time.Now().UTC() }
