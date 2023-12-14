package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("run echo.sh with envs", func(t *testing.T) {
		env := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}
		code := RunCmd([]string{"./testdata/echo.sh", "1", "2"}, env)
		require.Equal(t, code, OK)
	})
	t.Run("run echo.sh without env", func(t *testing.T) {
		env := Environment{}
		code := RunCmd([]string{"./testdata/echo.sh", "1", "2"}, env)
		require.Equal(t, code, OK)
	})
	t.Run("run echo.sh without args", func(t *testing.T) {
		env := Environment{}
		code := RunCmd([]string{"./testdata/echo.sh"}, env)
		require.Equal(t, code, OK)
	})
}
