package app

import (
	"context"
	"fmt"
	"log/slog"
	"notification-service/internal/clients"
	"notification-service/internal/transport/kafka"
	"notification-service/internal/transport/tg"
	"notification-service/pkg/closer"
	"notification-service/pkg/config"
	"notification-service/pkg/logger"
	"os"
	"os/signal"
	"strings"
	"time"
)

type App struct {
	log           *slog.Logger
	closer        *closer.Closer
	authClient    *clients.AuthClient
	kafkaConsumer *kafka.Consumer
	bot           *tg.Bot
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return nil, fmt.Errorf("app.NewApp load config: %w", err)
	}
	logger.Setup(cfg.AppEnv)
	logs := logger.With("service", "notification-service")
	logs.Info("initializing layers",
		slog.String("env", cfg.AppEnv))
	ctx = logger.WithContext(ctx, logs)

	authClient, err := clients.NewAuthClient(cfg.AuthServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("app.NewApp auth client: %w", err)
	}
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	bot, err := tg.NewBot(cfg.Token, authClient)

	consumer := kafka.NewConsumer(brokers, cfg.KafkaTopic, cfg.KafkaGroup, authClient, bot)

	cl := closer.New()

	cl.Add(func(ctx context.Context) error {
		slog.Info("closing tg bot")
		bot.BotAPI.StopReceivingUpdates()
		return nil
	})

	cl.Add(func(_ context.Context) error {
		slog.Info("closing kafka consumer")
		return consumer.Close()
	})

	cl.Add(func(_ context.Context) error {
		slog.Info("closing authClient")
		return authClient.Close()
	})

	return &App{
		log:           logs,
		closer:        cl,
		authClient:    authClient,
		kafkaConsumer: consumer,
		bot:           bot,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error)
	ctx = logger.WithContext(ctx, a.log)

	go func() {
		if err := a.kafkaConsumer.Start(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		a.bot.Start(ctx)
	}()

	a.log.Info("starting app")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case err := <-errCh:
		a.log.Error("app.Run server startup failed",
			slog.String("error", err.Error()))
	case sig := <-quit:
		a.log.Error("app.Run server shutdown",
			slog.Any("signal", sig))
	}

	a.log.Info("shutting down servers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		a.log.Error("shutdown errors", "err", err)
	}

	fmt.Println("Server Stopped")

	return nil
}
