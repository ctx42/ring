// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"maps"
	"os"
	"slices"
	"strings"
)

// Environ defines an interface for managing environment variables.
type Environ interface {
	// EnvLookup retrieves the value of the variable named by the key. If the
	// variable is present in the environment, the value (which may be empty)
	// is returned and the boolean is true. Otherwise, the returned value will
	// be empty and the boolean will be false.
	EnvLookup(key string) (string, bool)

	// EnvGet retrieves the value of the variable named by the key. It returns
	// the value, which will be empty if the variable is not present. To
	// distinguish between an empty value and an unset value, use
	// [Environ.EnvLookup].
	EnvGet(key string) string

	// EnvSet sets variable.
	EnvSet(key, value string)

	// EnvUnset unsets a single environment variable.
	EnvUnset(key string)

	// EnvAll returns environment as a slice of "key=value" entries. It
	// returns nil when the environment is empty.
	EnvAll() []string
}

var _ Environ = &Env{} // Compile time check.

// Env implements [Environ], storing environment variables.
type Env struct{ env map[string]string }

// NewEnv creates a new [Env] initialized with the given environment variables.
// If env is nil, an empty map is allocated. The input slice should contain
// "key=value" strings, as produced by [os.Environ].
func NewEnv(env []string) *Env {
	if env == nil {
		return &Env{env: make(map[string]string, 20)}
	}
	return &Env{env: EnvSplit(env)}
}

// EnvLookup retrieves the value of the environment variable named by the key
// from the given env slice. Returns the value (which may be empty) and true if
// the variable exists, or an empty string and false if it does not.
func (env *Env) EnvLookup(key string) (string, bool) {
	val, exist := env.env[key]
	return val, exist
}

// EnvGet retrieves the value of the environment variable named by the key from
// the given env slice. Returns the value or an empty string if not set. To
// distinguish between an empty value and an unset value, use [Env.EnvLookup].
func (env *Env) EnvGet(key string) string {
	val, _ := env.env[key]
	return val
}

// EnvSet sets the environment variable named by the key to the given value.
func (env *Env) EnvSet(key, value string) { env.env[key] = value }

// EnvSetFrom sets multiple environment variables from the given map.
// Overwrites existing variables with the same key.
func (env *Env) EnvSetFrom(src map[string]string) {
	for key, value := range src {
		env.EnvSet(key, value)
	}
}

// EnvSetWith sets multiple environment variables from the given slice. The
// input slice should contain "key=value" strings, as produced by [os.Environ].
// Overwrites existing variables with the same key.
func (env *Env) EnvSetWith(src []string) {
	env.EnvSetFrom(EnvSplit(src))
}

// EnvUnset unsets a single environment variable.
func (env *Env) EnvUnset(key string) { delete(env.env, key) }

// EnvAll returns environment as a slice of "key=value" entries. It returns nil
// when the environment is empty.
func (env *Env) EnvAll() []string {
	if len(env.env) == 0 {
		return nil
	}
	ret := make([]string, 0, len(env.env))
	for key, value := range env.env {
		ret = append(ret, key+"="+value)
	}
	return ret
}

// EnvClone returns a clone of the environment.
func (env *Env) EnvClone() *Env { return &Env{env: maps.Clone(env.env)} }

// EnvLookup retrieves the value of the "env" variable named by the key. If the
// variable is present in the "env", the value (which may be empty) is returned
// and the boolean is true. Otherwise, the returned value will be empty and the
// boolean will be false.
func EnvLookup(env []string, key string) (string, bool) {
	return NewEnv(env).EnvLookup(key)
}

// EnvGet retrieves the value of the "env" variable named by the key. It
// returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use [EnvLookup].
func EnvGet(env []string, key string) string {
	return NewEnv(env).EnvGet(key)
}

// EnvGetDefault retrieves the value of the "env" variable named by the key. If
// the key is not present in the environment, it will return def value.
func EnvGetDefault(env []string, key, def string) string {
	if val, exist := EnvLookup(env, key); exist {
		return val
	}
	return def
}

// EnvSet sets a single environment variable. Returns the modified slice.
func EnvSet(env []string, key, val string) []string {
	m := NewEnv(env)
	m.EnvSet(key, val)
	env = env[:0]
	for _, v := range m.EnvAll() {
		env = append(env, v)
	}
	return env
}

// EnvUnset unsets a single environment variable. Returns the modified slice.
func EnvUnset(env []string, key string) []string {
	m := NewEnv(env)
	m.EnvUnset(key)
	env = env[:0]
	for _, val := range m.EnvAll() {
		env = append(env, val)
	}
	return env
}

// EnvSplit parses [os.Environ] results and returns it as a key value map.
func EnvSplit(env []string) map[string]string {
	m := make(map[string]string, 10)
	for _, s := range env {
		if s == "" {
			continue
		}
		parts := strings.SplitN(s, "=", 2)
		if len(parts) == 2 {
			if parts[0] == "" {
				continue
			}
			m[parts[0]] = parts[1]
		}
	}
	return m
}

// EnvOrOs returns "env" if it's not nil, otherwise returns [os.Environ].
func EnvOrOs(env []string) []string {
	if env == nil {
		return os.Environ()
	}
	return env
}

// SetFrom sets environment variables from src map. Always returns a new slice.
func SetFrom(env []string, src map[string]string) []string {
	if len(src) == 0 {
		return slices.Clone(env)
	}
	ret := NewEnv(env)
	ret.EnvSetFrom(src)
	return ret.EnvAll()
}
