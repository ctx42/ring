// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package ring

import (
	"os"
	"slices"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

// envLookupTests are tabular tests for [Environ.EnvLookup] and [EnvLookup].
var envLookupTests = []struct {
	testN string

	env        []string
	findKey    string
	wantValue  string
	wantExists bool
}{
	{"found", []string{"key0=val0", "key1=val1"}, "key1", "val1", true},
	{"not found", []string{"key0=val0", "key1=val1"}, "key9", "", false},
	{"partial", []string{"key0=val0", "key1=val1"}, "key", "", false},
	{"empty env", []string{}, "key", "", false},
	{"empty key", []string{"key0=val0", "key1=val1"}, "", "", false},
	{
		"last value counts",
		[]string{"key0=val0", "key1=val1", "key0=abc"},
		"key0",
		"abc",
		true,
	},
}

// envGetTests are tabular tests for [Environ.EnvGet] and [EnvGet].
var envGetTests = []struct {
	testN string

	env       []string
	findKey   string
	wantValue string
}{
	{"found", []string{"key0=val0", "key1=val1"}, "key1", "val1"},
	{"not found", []string{"key0=val0", "key1=val1"}, "key9", ""},
	{"partial", []string{"key0=val0", "key1=val1"}, "key", ""},
	{"empty env", []string{}, "key", ""},
	{"empty key", []string{"key0=val0", "key1=val1"}, "", ""},
	{
		"last value counts",
		[]string{"key0=val0", "key1=val1", "key0=abc"},
		"key0",
		"abc",
	},
}

// envUnsetTests are tabular tests for [Environ.EnvUnset] and [EnvUnset].
var envUnsetTests = []struct {
	testN string

	env       []string
	deleteKey string
	wantEnv   []string
}{
	{"empty", nil, "A", nil},
	{
		"delete the first",
		[]string{"A=1", "B=2", "C=3"},
		"A",
		[]string{"B=2", "C=3"},
	},
	{
		"delete middle",
		[]string{"A=1", "B=2", "C=3"},
		"B",
		[]string{"A=1", "C=3"},
	},
	{
		"delete the last",
		[]string{"A=1", "B=2", "C=3"},
		"C",
		[]string{"A=1", "B=2"},
	},
}

func Test_NewEnv(t *testing.T) {
	t.Run("nil argument", func(t *testing.T) {
		// --- When ---
		have := NewEnv(nil)

		// --- Then ---
		assert.NotNil(t, have.env)
		assert.Len(t, 0, have.env)
	})

	t.Run("not nil argument", func(t *testing.T) {
		// --- Given ---
		want := []string{"A=1", "B=2"}

		// --- When ---
		have := NewEnv(want)

		// --- Then ---
		assert.Len(t, 2, have.env)
		assert.HasKeyValue(t, "A", "1", have.env)
		assert.HasKeyValue(t, "B", "2", have.env)
	})
}

func Test_Env_EnvLookup_tabular(t *testing.T) {
	for _, tc := range envLookupTests {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			env := NewEnv(tc.env)

			// --- When ---
			haveValue, haveExists := env.EnvLookup(tc.findKey)

			// --- Then ---
			assert.Equal(t, tc.wantValue, haveValue)
			assert.Equal(t, tc.wantExists, haveExists)
		})
	}
}

func Test_Env_EnvGet_tabular(t *testing.T) {
	for _, tc := range envGetTests {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			env := NewEnv(tc.env)

			// --- When ---
			haveValue := env.EnvGet(tc.findKey)

			// --- Then ---
			assert.Equal(t, tc.wantValue, haveValue)
		})
	}
}

func Test_Env_EnvSet(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{}}

		// --- When ---
		env.EnvSet("A", "1")

		// --- Then ---
		assert.Equal(t, map[string]string{"A": "1"}, env.env)
	})

	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{"A": "1", "B": "2"}}

		// --- When ---
		env.EnvSet("A", "2")

		// --- Then ---
		assert.Equal(t, map[string]string{"A": "2", "B": "2"}, env.env)
	})
}

func Test_Env_EnvSetFrom(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{"A": "1", "B": "2"}}

		// --- When ---
		env.EnvSetFrom(map[string]string{"A": "-1", "C": "3"})

		// --- Then ---
		want := map[string]string{"A": "-1", "B": "2", "C": "3"}
		assert.Equal(t, want, env.env)
	})

	t.Run("nil map", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{"A": "1", "B": "2"}}

		// --- When ---
		env.EnvSetFrom(nil)

		// --- Then ---
		assert.Equal(t, map[string]string{"A": "1", "B": "2"}, env.env)
	})
}

func Test_Env_EnvSetWith(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{"A": "1", "B": "2"}}

		// --- When ---
		env.EnvSetWith([]string{"A=-1", "C=3"})

		// --- Then ---
		want := map[string]string{"A": "-1", "B": "2", "C": "3"}
		assert.Equal(t, want, env.env)
	})

	t.Run("nil map", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{"A": "1", "B": "2"}}

		// --- When ---
		env.EnvSetWith(nil)

		// --- Then ---
		assert.Equal(t, map[string]string{"A": "1", "B": "2"}, env.env)
	})
}

func Test_Env_EnvUnset_tabular(t *testing.T) {
	for _, tc := range envUnsetTests {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			env := NewEnv(tc.env)

			// --- When ---
			env.EnvUnset(tc.deleteKey)

			// --- Then ---
			assert.Equal(t, tc.wantEnv, Sort(env.EnvAll()))
		})
	}
}

func Test_Env_EnvAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		env := &Env{env: map[string]string{"A": "1", "B": "2"}}

		// --- When ---
		have := env.EnvAll()

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(have))
	})

	t.Run("empty environment", func(t *testing.T) {
		// --- Given ---
		env := &Env{}

		// --- When ---
		have := env.EnvAll()

		// --- Then ---
		assert.Nil(t, have)
		assert.Len(t, 0, have)
	})
}

func Test_Env_EnvClone(t *testing.T) {
	// --- Given ---
	env := &Env{env: map[string]string{"A": "1"}}

	// --- When ---
	have := env.EnvClone()

	// --- Then ---
	assert.Equal(t, map[string]string{"A": "1"}, env.env)
	assert.NotSame(t, env.env, have.env)
	assert.NotSame(t, env, have)
	assert.Fields(t, 1, Env{})
}

func Test_EnvSet(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1"}

		// --- When ---
		have := EnvSet(env, "B", "2")

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2"}, Sort(have))
	})

	t.Run("existing", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1", "B=2", "C=3"}

		// --- When ---
		have := EnvSet(env, "B", "4")

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=4", "C=3"}, Sort(have))
	})
}

func Test_EnvUnset_tabular(t *testing.T) {
	for _, tc := range envUnsetTests {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			env := slices.Clone(tc.env)

			// --- When ---
			have := EnvUnset(env, tc.deleteKey)

			// --- Then ---
			assert.Equal(t, tc.wantEnv, Sort(have))
		})
	}
}

func Test_EnvLookup_tabular(t *testing.T) {
	for _, tc := range envLookupTests {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveValue, haveExists := EnvLookup(tc.env, tc.findKey)

			// --- Then ---
			assert.Equal(t, tc.wantValue, haveValue)
			assert.Equal(t, tc.wantExists, haveExists)
		})
	}
}

func Test_EnvGet_tabular(t *testing.T) {
	for _, tc := range envGetTests {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveValue := EnvGet(tc.env, tc.findKey)

			// --- Then ---
			assert.Equal(t, tc.wantValue, haveValue)
		})
	}
}

func Test_EnvGetDefault_tabular(t *testing.T) {
	tt := []struct {
		testN string

		env  []string
		key  string
		def  string
		want string
	}{
		{"get existing value", []string{"A=1", "B=2"}, "A", "x", "1"},
		{"get empty value", []string{"A=", "B=2"}, "A", "x", ""},
		{"get default value", []string{"A=1", "B=2"}, "C", "x", "x"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := EnvGetDefault(tc.env, tc.key, tc.def)

			// --- Ten ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_EnvSplit_tabular(t *testing.T) {
	tt := []struct {
		testN string

		env  []string
		want map[string]string
	}{
		{"empty slice", []string{}, map[string]string{}},
		{"empty entry", []string{""}, map[string]string{}},
		{"regular", []string{"A=B"}, map[string]string{"A": "B"}},
		{"multi equal sign", []string{"A=B=C"}, map[string]string{"A": "B=C"}},
		{"empty value", []string{"A="}, map[string]string{"A": ""}},
		{"just a key", []string{"A"}, map[string]string{}},
		{"just a value", []string{"=A"}, map[string]string{}},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := EnvSplit(tc.env)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}
func Test_EnvOrOs(t *testing.T) {
	t.Run("return env", func(t *testing.T) {
		// --- Given ---
		env := []string{"A=1"}

		// --- When ---
		have := EnvOrOs(env)

		// --- Then ---
		assert.Same(t, env, have)
	})

	t.Run("return os", func(t *testing.T) {
		// --- When ---
		have := EnvOrOs(nil)

		// --- Then ---
		assert.Equal(t, os.Environ(), have)
	})
}

func Test_SetFrom(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		// --- Given ---
		env := make([]string, 0, 10)
		env = append(env, "A=1", "B=2")

		// --- When ---
		have := SetFrom(env, map[string]string{"A": "-1", "C": "3"})

		// --- Then ---
		assert.Equal(t, []string{"A=-1", "B=2", "C=3"}, Sort(have))
		assert.NotSame(t, env, have)
	})

	t.Run("nil map", func(t *testing.T) {
		// --- Given ---
		env := make([]string, 0, 10)
		env = append(env, "A=1", "B=2")

		// --- When ---
		have := SetFrom(env, nil)

		// --- Then ---
		assert.Equal(t, []string{"A=1", "B=2"}, have)
		assert.NotSame(t, env, have)
	})
}
