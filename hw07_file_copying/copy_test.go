package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name         string
		offset       int64
		limit        int64
		expectedFile string
	}{
		{
			name:         "Full file copy",
			offset:       0,
			limit:        0,
			expectedFile: "out_offset0_limit0.txt",
		},
		{
			name:         "Limit 10 bytes",
			offset:       0,
			limit:        10,
			expectedFile: "out_offset0_limit10.txt",
		},
		{
			name:         "Limit 1000 bytes",
			offset:       0,
			limit:        1000,
			expectedFile: "out_offset0_limit1000.txt",
		},
		{
			name:         "Limit 10000 bytes",
			offset:       0,
			limit:        10000,
			expectedFile: "out_offset0_limit10000.txt",
		},
		{
			name:         "Offset 100, Limit 1000",
			offset:       100,
			limit:        1000,
			expectedFile: "out_offset100_limit1000.txt",
		},
		{
			name:         "Offset 6000, Limit 1000",
			offset:       6000,
			limit:        1000,
			expectedFile: "out_offset6000_limit1000.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outFile := filepath.Join("/tmp", "copy-test-"+tc.name+".txt")
			defer os.Remove(outFile)

			err := Copy(filepath.Join("testdata", "input.txt"), outFile, tc.offset, tc.limit)
			assert.NoError(t, err, "Ошибка копирования")

			expectedFile := filepath.Join("testdata", tc.expectedFile)

			outContent, err := os.ReadFile(outFile)
			assert.NoError(t, err, "Не удалось прочитать выходной файл")

			expectedContent, err := os.ReadFile(expectedFile)
			assert.NoError(t, err, "Не удалось прочитать эталонный файл")

			assert.Equal(t, expectedContent, outContent, "Содержимое файла не совпадает с эталоном")
		})
	}
}

func TestCopyErrorCases(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		srcFile := filepath.Join("testdata", "input.txt")
		outFile := filepath.Join("/tmp", "copy-error-test.txt")
		defer os.Remove(outFile)

		// получаем значение offset > длина файла
		srcFileInfo, err := os.Stat(srcFile)
		require.NoError(t, err)
		offset := srcFileInfo.Size() + 100

		err = Copy(srcFile, outFile, offset, 0)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("unsupported file (empty file)", func(t *testing.T) {
		emptyFile := filepath.Join("/tmp", "empty-file.txt")
		outFile := filepath.Join("/tmp", "copy-error-test.txt")
		defer os.Remove(outFile)
		defer os.Remove(emptyFile)

		_, err := os.Create(emptyFile)
		assert.NoError(t, err, "Не удалось создать пустой файл")

		err = Copy(emptyFile, outFile, 0, 0)
		assert.Error(t, err, "Ожидалась ошибка при копировании пустого файла")
		assert.ErrorIs(t, err, ErrUnsupportedFile, "Ожидалась ошибка ErrUnsupportedFile для пустого файла")
	})

	t.Run("unsupported file (directory)", func(t *testing.T) {
		outFile := filepath.Join("/tmp", "copy-error-test.txt")
		defer os.Remove(outFile)

		err := Copy("testdata", outFile, 0, 0)
		assert.Error(t, err, "Ожидалась ошибка при копировании директории")
		assert.ErrorIs(t, err, ErrUnsupportedFile, "Ожидалась ошибка ErrUnsupportedFile для директории")
	})

	t.Run("not exist file", func(t *testing.T) {
		toFile := filepath.Join("/tmp", "copy-error-test.txt")
		defer os.Remove(toFile)
		err := Copy("unknown.txt", toFile, 0, 0)

		assert.Error(t, err, "Ожидалась ошибка при копировании несуществующего файла")
		assert.ErrorIs(t, err, os.ErrNotExist, "Ожидалась ошибка os.ErrNotExist")
	})
}
