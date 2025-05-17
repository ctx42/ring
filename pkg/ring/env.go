// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"os"
	"slices"
	"strings"
)

// Environ represents environment.
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

	// EnvGetAll returns environment as a slice of "key=value" entries. It
	// returns nil when the environment is empty.
	EnvGetAll() []string
}

var _ Environ = Env{} // Compile time check.

// Env represents environment.
type Env struct{ e map[string]string }

// NewEnv returns new instance of [Env] initialized with `env`. If `env` is nil,
// the new map will be allocated.
func NewEnv(env []string) Env {
	if env == nil {
		return Env{e: make(map[string]string, 20)}
	}
	return Env{e: EnvSplit(env)}
}

// EnvLookup retrieves the value of the variable named by the key. If the
// variable is present in the environment, the value (which may be empty) is
// returned and the boolean is true. Otherwise, the returned value will be
// empty and the boolean will be false.
func (env Env) EnvLookup(key string) (string, bool) {
	val, exist := env.e[key]
	return val, exist
}

// EnvGet retrieves the value of the variable named by the key. It returns the
// value, which will be empty if the variable is not present. To distinguish
// between an empty value and an unset value, use [Env.EnvLookup].
func (env Env) EnvGet(key string) string {
	val, _ := env.e[key]
	return val
}

// EnvGetAll returns environment as a slice of "key=value" entries. It returns
// nil when the environment is empty.
func (env Env) EnvGetAll() []string {
	var list []string
	for key, value := range env.e {
		list = append(list, key+"="+value)
	}
	return list
}

// EnvSet sets variable.
func (env Env) EnvSet(key, value string) { env.e[key] = value }

// EnvSetFrom sets environment variables from given map.
func (env Env) EnvSetFrom(src map[string]string) {
	for key, value := range src {
		env.EnvSet(key, value)
	}
}

// EnvUnset unsets a single environment variable.
func (env Env) EnvUnset(key string) { delete(env.e, key) }

// EnvIsNil returns true if no memory was allocated for environment.
func (env Env) EnvIsNil() bool { return env.e == nil }

// EnvLookup retrieves the value of the `env` variable named by the key. If the
// variable is present in the `env` the value (which may be empty) is returned
// and the boolean is true. Otherwise, the returned value will be empty and the
// boolean will be false.
func EnvLookup(env []string, key string) (string, bool) {
	return NewEnv(env).EnvLookup(key)
}

// EnvGet retrieves the value of the `env` variable named by the key. It
// returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use [EnvLookup].
func EnvGet(env []string, key string) string {
	return NewEnv(env).EnvGet(key)
}

// EnvGetOr retrieves the value of the `env` variable named by the key. If the
// key is not present in the environment, it will return def value.
func EnvGetOr(env []string, key, def string) string {
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
	for _, val := range m.EnvGetAll() {
		env = append(env, val)
	}
	return env
}

// EnvUnset unsets a single environment variable. Returns the modified slice.
func EnvUnset(env []string, key string) []string {
	m := NewEnv(env)
	m.EnvUnset(key)
	env = env[:0]
	for _, val := range m.EnvGetAll() {
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

// EnvOrOs returns `env` if it's not nil, otherwise returns [os.Environ].
func EnvOrOs(env []string) []string {
	if env == nil {
		return os.Environ()
	}
	return env
}

// SetFrom sets environment variables from src map. Always returns new slice.
func SetFrom(env []string, src map[string]string) []string {
	if len(src) == 0 {
		return slices.Clone(env)
	}
	ret := NewEnv(env)
	ret.EnvSetFrom(src)
	return ret.EnvGetAll()
}
