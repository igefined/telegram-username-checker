package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"go.uber.org/zap"

	"github.com/igefined/go-kit/config"
	"github.com/igefined/go-kit/log"
)

func main() {
	var (
		ctx       = config.SigTermIntCtx()
		logger, _ = log.NewLogger(zap.DebugLevel)
		cfg       = Config{}
	)

	if err := config.GetConfig("", &cfg, envs); err != nil {
		logger.Fatal("unable to load config", zap.Error(err))
	}

	file, err := os.OpenFile(fmt.Sprintf("./files/%s", cfg.UsernameList), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Fatal("unable to open config file", zap.Error(err))
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var outputUsernames = make([]string, 0, len(lines))
	client, err := telegram.ClientFromEnvironment(telegram.Options{Logger: logger.Logger})
	if err != nil {
		logger.Fatal("unable to create telegram client", zap.Error(err))
	}

	if err = client.Run(ctx, func(ctx context.Context) error {
		stop, err := bg.Connect(client)
		if err != nil {
			return err
		}
		defer func() { _ = stop() }()

		if _, err = client.Auth().Status(ctx); err != nil {
			return err
		}

		var (
			api      = client.API()
			isExists bool
		)

		for _, line := range lines {
			logger.Info("account check username", zap.String("username", line))

			isExists, err = api.AccountCheckUsername(ctx, line)
			if err != nil {
				return err
			}

			if !isExists {
				outputUsernames = append(outputUsernames, line)
			}
		}

		return nil
	}); err != nil {
		logger.Fatal("telegram run application error", zap.Error(err))
	}
}
