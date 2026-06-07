// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring_test

import (
	"bytes"
	"fmt"
	"time"

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

// ExampleWithClock shows how to inject a fixed clock so that
// time-dependent code produces deterministic output in tests.
func ExampleWithClock() {
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	rng := ring.New(ring.WithClock(func() time.Time { return fixed }))

	fmt.Println(rng.Clock()().Format(time.DateOnly))
	// Output:
	// 2024-01-15
}

// ExampleRing_MetaSet shows how to store and retrieve arbitrary metadata
// on a Ring, useful for passing context (e.g. a trace ID or config) from
// a top-level command down to subcommand handlers without extra parameters.
func ExampleRing_MetaSet() {
	rng := ring.New()
	rng.MetaSet("trace-id", "abc-123")

	fmt.Println(rng.MetaGet("trace-id"))
	// Output:
	// abc-123
}

// ExampleRing_Clone shows how to create an independent subcommand context
// from a parent Ring. The clone gets its own I/O and environment, but
// shares the parent's metadata map — changes on either side are visible
// to both.
func ExampleRing_Clone() {
	parent := ring.New()
	parent.MetaSet("trace-id", "xyz-789")

	child := parent.Clone()
	child.SetArgs([]string{"--verbose"})

	fmt.Println(child.Args())
	fmt.Println(child.MetaGet("trace-id"))
	// Output:
	// [--verbose]
	// xyz-789
}
