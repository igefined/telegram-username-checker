package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gotd/td/telegram"

	"github.com/igefined/go-kit/log"
)

const outputDir = "./out"

type App struct {
	logger *log.Logger
	cfg    *Config

	client *telegram.Client
}

func NewApp(logger *log.Logger, cfg *Config) (*App, error) {
	client, err := telegram.ClientFromEnvironment(telegram.Options{Logger: logger.Logger})
	if err != nil {
		return nil, err
	}

	return &App{
		logger: logger,
		cfg:    cfg,
		client: client,
	}, nil
}

func (a *App) RunChecker(ctx context.Context, usernames []string) error {
	var out = make([]string, 0)

	if err := a.client.Run(ctx, func(ctx context.Context) error {
		for i, username := range usernames {
			if i%2 == 0 {
				out = append(out, username)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	_ = os.Mkdir(outputDir, 0755)

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", outputDir, a.cfg.OutputPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	for _, username := range out {
		if _, err = file.WriteString(username + "\n"); err != nil {
			return err
		}
	}

	return err
}
