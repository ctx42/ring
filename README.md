[![Go Report Card](https://goreportcard.com/badge/github.com/ctx42/ring)](https://goreportcard.com/report/github.com/ctx42/ring)
[![GoDoc](https://img.shields.io/badge/api-Godoc-blue.svg)](https://pkg.go.dev/github.com/ctx42/ring)
![Tests](https://github.com/ctx42/ring/actions/workflows/go.yml/badge.svg?branch=master)

<!-- TOC -->
* [The `ring` Package](#the-ring-package)
* [Key Features](#key-features)
* [Installation](#installation)
* [Usage](#usage)
  * [Production Code](#production-code)
  * [Test Code](#test-code)
  * [Metadata](#metadata)
  * [Subcommand Context](#subcommand-context)
  * [Deterministic Time](#deterministic-time)
<!-- TOC -->

# The `ring` Package

The `ring` package provides utilities to manage a program's execution context,
enabling dependency injection for standard I/O streams, environment variables,
arguments, and time. By avoiding reliance on global state (e.g., `os.Stdin`, 
`os.Environ`), it simplifies testing and improves code modularity.

![ring.png](doc/ring.png)

# Key Features

- **Dependency Injection**: inject custom standard I/O streams, environment variables, program name, arguments, and a clock.
- **Test-Friendly**: simplifies mocking of global dependencies for unit tests.
- **Metadata Support**: store and manage arbitrary key-value metadata.

# Installation

To use `ring` in your project, add it as a dependency:

```shell
go get github.com/ctx42/ring
```

# Usage

The `ring` package centers around the `Ring` type, which encapsulates a 
program's execution context. Below are examples demonstrating its core 
functionality.

## Production Code

Production code example.

```go
package main

import (
    "context"
    "os"

    "github.com/ctx42/ring/pkg/ring"

    "github.com/user/project/cmd"
)

func main() {
    // Default Ring:
    //  - Standard I/O: [os.Stdin], [os.Stdout], [os.Stderr]
    //  - Environment: [os.Environ]
    //  - Clock: [NowUTC]
    //  - Args: os.Args[1:]
    //  - Name: os.Args[0]
    //  - Metadata: empty map
    //  - Filesystem: none
    rng := ring.New(ring.WithFS(os.DirFS("some/path")))

    ctx := context.Background()
    
    exitCode := cmd.Main(ctx, rng) // Call the real application entrypoint.
    os.Exit(exitCode)
}
```

This way the `cmd.Main` becomes really easy to test. 

## Test Code

Use a buffer for stdout to capture and verify program output without
touching `os.Stdout`:

<!-- gmdoceg:ExampleNew_inTest -->
```go
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
```

## Metadata

`Ring` carries an arbitrary `map[string]any` that subcommand handlers
can read and write without extra function parameters:

<!-- gmdoceg:ExampleRing_MetaSet -->
```go
rng := ring.New()
rng.MetaSet("trace-id", "abc-123")

fmt.Println(rng.MetaGet("trace-id"))
// Output:
// abc-123
```

## Subcommand Context

`Clone` produces an independent `Ring` with its own I/O and environment.
The metadata map is shared intentionally — a parent command can set a
trace ID or loaded config once and every clone sees it:

<!-- gmdoceg:ExampleRing_Clone -->
```go
parent := ring.New()
parent.MetaSet("trace-id", "xyz-789")

child := parent.Clone()
child.SetArgs([]string{"--verbose"})

fmt.Println(child.Args())
fmt.Println(child.MetaGet("trace-id"))
// Output:
// [--verbose]
// xyz-789
```

## Deterministic Time

Inject a fixed clock to make time-dependent code produce stable output
in tests:

<!-- gmdoceg:ExampleWithClock -->
```go
fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
rng := ring.New(ring.WithClock(func() time.Time { return fixed }))

fmt.Println(rng.Clock()().Format(time.DateOnly))
// Output:
// 2024-01-15
```
