package main

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/igefined/go-kit/config"
	"github.com/igefined/go-kit/log"
)

const usernameFilesDir = "./files"

func main() {
	var (
		ctx       = config.SigTermIntCtx()
		logger, _ = log.NewLogger(zap.DebugLevel)
		cfg       = Config{}
	)

	if err := config.GetConfig("", &cfg, envs); err != nil {
		logger.Fatal("unable to load config", zap.Error(err))
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", usernameFilesDir, cfg.UsernameList), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Fatal("unable to open config file", zap.Error(err))
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	app, err := NewApp(logger, &cfg)
	if err != nil {
		logger.Fatal("unable to create app", zap.Error(err))
	}

	if err = app.RunChecker(ctx, lines); err != nil {
		logger.Fatal("unable to create app", zap.Error(err))
	}
}
