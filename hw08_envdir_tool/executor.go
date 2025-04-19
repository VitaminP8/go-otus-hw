package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/VitaminP8/go-otus-hw/hw08_envdir_tool/logger"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	logger.Logger.Info("running command: ", "cmd", cmd)
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	// записываем нынешние переменные окружения
	currEnv := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			currEnv[parts[0]] = parts[1]
		}
	}

	// обновляем переменные окружения (из envDir)
	for key, value := range env {
		if value.NeedRemove {
			delete(currEnv, key)
		} else {
			currEnv[key] = value.Value
		}
	}

	// преобразуем словарь обратно в формат KEY+VALUE
	finalEnv := make([]string, 0, len(currEnv))
	for key, value := range currEnv {
		finalEnv = append(finalEnv, key+"="+value)
	}
	logger.Logger.Debug("final env", "env", finalEnv)

	command.Env = finalEnv

	// Запускаем команду
	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			logger.Logger.Warn("Command exited with error", "code", exitErr.ExitCode())
			return exitErr.ExitCode()
		}
		logger.Logger.Error("command failed to run", "err", err)
		return 1
	}

	logger.Logger.Info("Command completed successfully")
	return command.ProcessState.ExitCode()
}
