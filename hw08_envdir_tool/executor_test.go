package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Correct run with env", func(t *testing.T) {
		cmd := []string{"bash", "-c", "echo VAR=$VAR"}
		env := Environment{
			"VAR": {Value: "123", NeedRemove: false},
		}

		// временно перенаправляем stdout
		orig := os.Stdout
		r, w, err := os.Pipe()
		require.NoError(t, err)
		os.Stdout = w

		code := RunCmd(cmd, env)

		// закрываем и читаем результат
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = orig

		require.Equal(t, 0, code)
		require.Equal(t, "VAR=123", strings.TrimSpace(string(out)))
	})

	t.Run("Should remove env var", func(t *testing.T) {
		os.Setenv("REMOVE_ME", "dont care")

		cmd := []string{"bash", "-c", "echo REMOVE_ME=$REMOVE_ME"}
		env := Environment{
			"REMOVE_ME": {Value: "", NeedRemove: true},
		}

		orig := os.Stdout
		r, w, err := os.Pipe()
		require.NoError(t, err)
		os.Stdout = w

		code := RunCmd(cmd, env)

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = orig

		require.Equal(t, 0, code)
		require.Equal(t, "REMOVE_ME=", strings.TrimSpace(string(out)))
	})

	t.Run("Return non-zero code on bad command", func(t *testing.T) {
		cmd := []string{"bash", "-c", "exit 42"}
		env := Environment{}
		code := RunCmd(cmd, env)

		require.Equal(t, 42, code)
	})

	t.Run("Return 1 on exec error", func(t *testing.T) {
		cmd := []string{"/nonexistent"}
		env := Environment{}
		code := RunCmd(cmd, env)

		require.Equal(t, 1, code)
	})
}
