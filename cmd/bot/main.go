package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mauzec/ibot-things/config"
	"github.com/mauzec/ibot-things/internal/bot"
	"github.com/mauzec/ibot-things/internal/data"
	"github.com/mauzec/ibot-things/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapAdapter struct {
	log *zap.SugaredLogger
}

func (z *zapAdapter) Printf(format string, v ...any) {
	z.log.Infof(format, v...)
}

func (z *zapAdapter) Println(v ...any) {
	z.log.Info(v...)
}

func main() {
	if err := logger.Init(true); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	zapLogger := logger.Logger().Desugar()
	zap.RedirectStdLog(zapLogger)
	zap.RedirectStdLogAt(zapLogger, zapcore.DebugLevel)

	tgbotapi.SetLogger(&zapAdapter{log: logger.Logger()})

	cfg, err := config.LoadConfig("app", "env", "config")
	if err != nil {
		logger.Logger().Fatalf("failed to load config: %v", err)
	}

	if cfg.DeviceAddr == "" || cfg.DevicePort == 0 || cfg.PollIntervalMs <= 0 {
		logger.Logger().Fatalf("device address or poll interval is not set")
	}

	interval := time.Duration(cfg.PollIntervalMs) * time.Millisecond

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	esp := data.NewEspProvider(
		cfg.DeviceAddr,
		cfg.DevicePort,
		cfg.DataPath,
		interval,
	)
	esp.Go(ctx)

	b, err := bot.NewBotFromConfig(&cfg)
	if err != nil {
		logger.Logger().Fatalf("failed to create bot: %v", err)
	}

	b.Run(ctx)
}
