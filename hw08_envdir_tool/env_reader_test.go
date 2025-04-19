package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Valid env", func(t *testing.T) {
		dir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(dir, "BAR"), []byte("bar\n123"), 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "EMPTY"), []byte(""), 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "TRIM"), []byte("trim \t\n "), 0o644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "NULLBYTE"), []byte("value\x00123"), 0o644))

		env, err := ReadDir(dir)
		require.NoError(t, err)

		require.Equal(t, 4, len(env))

		require.Equal(t, "bar", env["BAR"].Value)
		require.False(t, env["BAR"].NeedRemove)

		require.Equal(t, "", env["EMPTY"].Value)
		require.True(t, env["EMPTY"].NeedRemove)

		require.Equal(t, "trim", env["TRIM"].Value)
		require.False(t, env["TRIM"].NeedRemove)

		require.Equal(t, "value\n123", env["NULLBYTE"].Value)
		require.False(t, env["NULLBYTE"].NeedRemove)
	})

	t.Run("Invalid env", func(t *testing.T) {
		dir := t.TempDir()

		require.NoError(t, os.WriteFile(filepath.Join(dir, "BAD=FILENAME"), []byte("dont care"), 0o644))

		env, err := ReadDir(dir)
		require.ErrorContains(t, err, "invalid filename")
		require.Nil(t, env)
	})
}
