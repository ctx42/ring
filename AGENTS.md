This file provides guidance to AI agents (and human contributors) when working
with code in this repository.

## Commands

```sh
# Run all tests with race detector
go test -race ./...

# Run tests for a single package
go test -race ./pkg/ring/...

# Run a single test
go test -race -run TestName ./pkg/ring/...

# Tidy dependencies
go mod tidy
```

CI runs `go test -v -race ./...` on every push/PR to master.

## Architecture

`github.com/ctx42/ring` is a small Go library providing a **program execution
context** (`Ring`) that bundles standard I/O streams, environment variables,
arguments, a clock, and arbitrary metadata into a single injectable value. It is
designed to be passed to CLI subcommands so they work against swappable streams
and env instead of global OS state.

### Package layout

- **`pkg/ring`** — core package. Three composable types are embedded into
  `Ring`:
    - `IO` — wraps `stdin`/`stdout`/`stderr` (`io.Reader`/`io.Writer`).
      Implements `Streamer`.
    - `Env` — wraps environment variables as a `map[string]string`. Implements
      `Environ`. Both `Env` and `IO` are embedded as unexported type aliases (
      `hidEnv`, `hidIO`) so their methods surface on `Ring` without exposing the
      field names.
    - `Ring` itself adds: program name, args `[]string`, `Clock` func, `fs.FS`,
      and a `map[string]any` metadata bag.
    - `Clone()` deep-copies everything except the metadata map (shared by
      reference across clones by design).

- **`pkg/ring/ringtest`** — test helper. `ringtest.Tester` wraps a `Ring` with
  `bytes.Buffer`/`iokit.Buffer` streams and exposes `WetStdout()`/
  `WetStderr()` (assert written) vs the default `DryBuffer` (assert not
  written). Call `tst.Ring(args...)` to get a ready-to-use `*ring.Ring` for a
  test case.

### Key design patterns

- All configuration goes through functional options (`Option func(*Ring)`)
  passed to `New(...)`.
- Sentinel errors (`ErrReqMeta`, `ErrInvMeta`, `ErrNoFsAccess`) are declared in
  `ring.go` and used for typed error checking by callers.
- `Ring` carries metadata (`MetaSet`/`MetaGet`/`MetaLookup`/`MetaDelete`/
  `MetaAll`) for passing context between CLI layers without extra function
  parameters.
- `Env` has both method-receiver API (`*Env`) and standalone helpers (
  `EnvLookup`, `EnvGet`, `EnvSet`, `EnvUnset`, `EnvSplit`) that operate on
  `[]string` slices — useful when working with raw `os.Environ()` output without
  constructing an `Env`.
