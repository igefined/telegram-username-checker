package main

import (
	"sync"

	"github.com/igefined/go-kit/config"
)

var envs = []*config.EnvVar{
	config.NewEnvVar(
		"file_username_list",
		"FILE_USERNAME_LIST",
		"",
		"Filepath to file with username list",
	),
	config.NewEnvVar(
		"file_output_path",
		"FILE_OUTPUT_PATH",
		"",
		"Filepath to result file",
	),
}

type (
	FilesConfig struct {
		UsernameList string `mapstructure:"file_username_list"`
		OutputPath   string `mapstructure:"file_output_path"`
	}

	Config struct {
		sync.RWMutex

		FilesConfig `mapstructure:",squash"`
	}
)
