package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/udhos/equalfile"
)

func TestCopy(t *testing.T) {
	cmp := equalfile.New(nil, equalfile.Options{})
	tests := []struct {
		name        string
		inputPath   string
		comparePath string
		offset      int64
		limit       int64
	}{
		{
			name:        "full copy",
			inputPath:   "testdata/input.txt",
			comparePath: "testdata/out_offset0_limit0.txt",
			offset:      int64(0),
			limit:       int64(0),
		},
		{
			name:        "copy with offset=0/limit=10",
			inputPath:   "testdata/input.txt",
			comparePath: "testdata/out_offset0_limit10.txt",
			offset:      int64(0),
			limit:       int64(10),
		},
		{
			name:        "copy with offset=0/limit=1000",
			inputPath:   "testdata/input.txt",
			comparePath: "testdata/out_offset0_limit1000.txt",
			offset:      int64(0),
			limit:       int64(1000),
		},
		{
			name:        "copy with offset=0/limit=10000",
			inputPath:   "testdata/input.txt",
			comparePath: "testdata/out_offset0_limit10000.txt",
			offset:      int64(0),
			limit:       int64(10000),
		},
		{
			name:        "copy with offset=100/limit=1000",
			inputPath:   "testdata/input.txt",
			comparePath: "testdata/out_offset100_limit1000.txt",
			offset:      int64(100),
			limit:       int64(1000),
		},
		{
			name:        "copy with offset=6000/limit=1000",
			inputPath:   "testdata/input.txt",
			comparePath: "testdata/out_offset6000_limit1000.txt",
			offset:      int64(6000),
			limit:       int64(1000),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp("", "tmpfile-")
			require.NoError(t, err)
			defer os.Remove(f.Name())
			outputPath := f.Name()
			err = Copy(tc.inputPath, outputPath, tc.offset, tc.limit)
			require.NoError(t, err)
			equal, err := cmp.CompareFile(outputPath, tc.comparePath)
			require.NoError(t, err)
			require.True(t, equal)
		})
	}

	t.Run("offset > file size", func(t *testing.T) {
		fileInfo, err := os.Stat("testdata/input.txt")
		require.NoError(t, err)
		f, err := os.CreateTemp("", "tmpfile-")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		outputPath := f.Name()
		err = Copy("testdata/input.txt", outputPath, fileInfo.Size()+1, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		f, err := os.CreateTemp("", "tmpfile-")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		var (
			inputPath  = "/dev/urandom"
			outputPath = f.Name()
			offset     = int64(0)
			limit      = int64(0)
		)
		err = Copy(inputPath, outputPath, offset, limit)
		require.Error(t, err)
		require.Equal(t, ErrUnsupportedFile, err)
	})
}
