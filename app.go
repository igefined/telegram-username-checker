package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"

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

	flow := auth.NewFlow(Terminal{}, auth.SendCodeOptions{})

	if err := a.client.Run(ctx, func(ctx context.Context) error {
		if err := a.client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		api := a.client.API()

		for _, username := range usernames {
			_, err := api.AccountCheckUsername(ctx, username)
			if err == nil {
				out = append(out, username)
				continue
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
