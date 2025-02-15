package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inPath := "testdata/input.txt"
	outPath := "out.txt"
	t.Run("offset exceed", func(t *testing.T) {
		err := Copy(inPath, outPath, 7000, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("offset negative", func(t *testing.T) {
		err := Copy(inPath, outPath, -1, 0)
		require.Equal(t, err, ErrNegativeValues)
	})

	t.Run("limit negative", func(t *testing.T) {
		err := Copy(inPath, outPath, 2, -1)
		require.Equal(t, err, ErrNegativeValues)
	})

	t.Run("successful case, offset 100, limit 10000", func(t *testing.T) {
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 100, 1000)
		require.NoError(t, err)

		stat, _ := os.Stat(outPath)
		require.Equal(t, int64(1000), stat.Size())
	})

	t.Run("successful case, offset 0, limit 0", func(t *testing.T) {
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 0, 0)
		require.NoError(t, err)

		sStat, err := os.Stat(inPath)
		require.NoError(t, err)

		oStat, err := os.Stat(outPath)
		require.NoError(t, err)

		require.Equal(t, sStat.Size(), oStat.Size())
	})

	t.Run("file doesn't exist", func(t *testing.T) {
		err := Copy("testdata/non_existent.txt", outPath, 0, 0)
		require.Error(t, err)
	})

	t.Run("limit is exceeded", func(t *testing.T) {
		defer os.Remove(outPath)

		err := Copy(inPath, outPath, 0, 999999)
		require.NoError(t, err)

		sStat, err := os.Stat(inPath)
		require.NoError(t, err)

		oStat, err := os.Stat(outPath)
		require.NoError(t, err)

		require.Equal(t, sStat.Size(), oStat.Size())
	})
}
