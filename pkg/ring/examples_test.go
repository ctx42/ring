// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring_test

import (
	"bytes"
	"fmt"

	"github.com/ctx42/ring/pkg/ring"
)

// ExampleNew_inTest shows the typical test pattern: inject a buffer for
// stdout and verify the program's output without touching os.Stdout.
func ExampleNew_inTest() {
	// greet simulates a CLI function that writes to the ring's stdout.
	greet := func(rng *ring.Ring) {
		name := rng.EnvGet("USER_NAME")
		fmt.Fprintf(rng.Stdout(), "Hello, %s!\n", name)
	}

	var sout bytes.Buffer
	rng := ring.New(
		ring.WithEnv([]string{"USER_NAME=Alice"}),
		ring.WithArgs([]string{"--verbose"}),
	)
	rng.SetStdout(&sout)

	greet(rng)

	fmt.Print(sout.String())
	// Output:
	// Hello, Alice!
}
