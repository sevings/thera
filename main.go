package main

import (
	"os"
	"os/signal"
	"time"

	"thera/internal/thera"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func main() {
	cfg, err := thera.LoadConfig()
	if err != nil {
		panic(err)
	}

	var zapLogger *zap.Logger
	if cfg.Release {
		zapLogger, err = zap.NewProduction(zap.WithCaller(false))
	} else {
		zapLogger, err = zap.NewDevelopment(zap.WithCaller(false))
	}
	if err != nil {
		panic(err)
	}
	defer func() { _ = zapLogger.Sync() }()

	zap.ReplaceGlobals(zapLogger)
	zap.RedirectStdLog(zapLogger)
	logger := zapLogger.Sugar()

	db, ok := thera.LoadDatabase(cfg.DBPath)
	if !ok {
		logger.Panic("can't load database")
	}

	th := thera.NewThera(db, cfg)
	bot := thera.NewBot(db)

	pref := tele.Settings{
		Token:   cfg.TgToken,
		Poller:  &tele.LongPoller{Timeout: 30 * time.Second},
		OnError: bot.LogError,
	}

	api, err := tele.NewBot(pref)
	if err != nil {
		logger.Panic(err)
	}

	bot.Start(cfg, api, th)
	defer bot.Stop()

	err = th.Start()
	if err != nil {
		logger.Panic(err.Error())
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
