package main

import (
	"os"

	"github.com/VitaminP8/go-otus-hw/hw08_envdir_tool/logger"
)

func main() {
	if len(os.Args) < 3 {
		logger.Logger.Error("Not enough arguments", "args", os.Args)
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2:]

	logger.Logger.Info("Starting envdir", "envDir", envDir, "command", command)

	// Читаем переменные окружения из директории
	env, err := ReadDir(envDir)
	if err != nil {
		logger.Logger.Error("Failed to read envdir", "err", err)
		os.Exit(1)
	}

	// Запускаем команду
	exitCode := RunCmd(command, env)
	// Завершаем процесс с тем же кодом
	logger.Logger.Info("Exiting with code", "code", exitCode)
	os.Exit(exitCode)
}
