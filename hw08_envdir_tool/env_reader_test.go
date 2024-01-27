package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("getting correct envs from testdata dir", func(t *testing.T) {
		env := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}
		result, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, result)
	})

	t.Run("invalid file name", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "tests")
		require.NoError(t, err)
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				require.NoError(t, err)
			}
		}(dir)
		_, err = os.Create(path.Join(dir, "tests=tests"))
		require.NoError(t, err)
		result, err := ReadDir(dir)
		require.Equal(t, err, ErrInvalidFileName)
		require.Empty(t, result)
	})

	t.Run("replacing 0x00 with \n", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "tests")
		require.NoError(t, err)
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				require.NoError(t, err)
			}
		}(dir)
		file, err := os.Create(path.Join(dir, "tests"))
		require.NoError(t, err)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				require.NoError(t, err)
			}
		}(file)
		_, err = file.WriteString("tests\x00data \x00")
		require.NoError(t, err)
		result, err := ReadDir(dir)
		require.NoError(t, err)
		env := Environment{"tests": EnvValue{
			Value:      "tests\ndata",
			NeedRemove: false,
		}}
		fmt.Println(result)
		require.Equal(t, env, result)
	})

	t.Run("no env files in dir", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "tests")
		require.NoError(t, err)
		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				require.NoError(t, err)
			}
		}(dir)
		result, err := ReadDir(dir)
		require.NoError(t, err)
		env := Environment{}
		require.Equal(t, env, result)
	})
}
