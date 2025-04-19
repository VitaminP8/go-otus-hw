package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/VitaminP8/go-otus-hw/hw08_envdir_tool/logger"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		if strings.Contains(fileName, "=") {
			logger.Logger.Warn("Invalid filename: contains '='", "file", fileName)
			return nil, errors.New("invalid filename: " + fileName)
		}

		filePath := filepath.Join(dir, fileName)
		file, err := os.Open(filePath)
		if err != nil {
			logger.Logger.Error("Failed to open file", "file", filePath, "err", err)
			return nil, err
		}
		defer file.Close()

		fileStat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if fileStat.Size() == 0 {
			logger.Logger.Debug("Empty file", "file", fileName)
			env[fileName] = EnvValue{"", true}
			continue
		}

		// читаем первую строку
		scanner := bufio.NewScanner(file)
		var line string
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, " \t")
		line = strings.ReplaceAll(line, "\x00", "\n")

		logger.Logger.Debug("Reading file", "file", fileName, "line", line)
		env[fileName] = EnvValue{line, false}
	}

	return env, nil
}
