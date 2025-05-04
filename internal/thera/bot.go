package thera

import (
	"context"
	"fmt"
	"strings"
	"thera/internal/letta"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	api   BotAPI
	db    *DB
	log   *zap.SugaredLogger
	cfg   Config
	thera *Thera
}

type BotAPI interface {
	ChatByID(id int64) (*tele.Chat, error)
	Send(to tele.Recipient, what any, opts ...any) (*tele.Message, error)
	Handle(endpoint any, h tele.HandlerFunc, m ...tele.MiddlewareFunc)
	Use(middlewares ...tele.MiddlewareFunc)
	Start()
	Stop()
}

func NewBot(db *DB) *Bot {
	bot := &Bot{
		db:  db,
		log: zap.L().Named("bot").Sugar(),
	}

	return bot
}

func (bot *Bot) Logger() *zap.SugaredLogger {
	return bot.log
}

func (bot *Bot) Start(cfg Config, api BotAPI, thera *Thera) {
	bot.api = api
	bot.cfg = cfg
	bot.thera = thera

	bot.api.Use(middleware.Recover())
	bot.api.Use(bot.logMessage)

	bot.api.Handle(tele.OnText, bot.handleText)
	bot.api.Handle("/start", bot.handleStart)

	go func() {
		bot.log.Info("starting bot")
		bot.api.Start()
		bot.log.Info("bot stopped")
	}()
}

func (bot *Bot) Stop() {
	bot.api.Stop()
}

func (bot *Bot) logMessage(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		beginTime := time.Now().UnixNano()

		err := next(c)

		endTime := time.Now().UnixNano()
		duration := float64(endTime-beginTime) / 1000000

		isCmd := len(c.Text()) > 0 && c.Text()[0] == '/' && len(c.Entities()) == 1
		isAction := c.Callback() != nil
		var action string
		if isCmd {
			action = c.Text()
		} else if isAction {
			action = strings.TrimSpace(strings.Split(c.Callback().Data, "|")[0])
		}
		bot.log.Infow("user message",
			"chat_id", c.Chat().ID,
			"chat_type", c.Chat().Type,
			"user_id", c.Sender().ID,
			"user_name", c.Sender().Username,
			"is_cmd", isCmd,
			"is_action", isAction,
			"action", action,
			"size", len(c.Text()),
			"dur", fmt.Sprintf("%.2f", duration),
			"err", err)

		return err
	}
}

func (bot *Bot) LogError(err error, c tele.Context) {
	if c == nil {
		bot.log.Errorw("error", "err", err)
	} else {
		isCmd := len(c.Text()) > 0 && c.Text()[0] == '/' && len(c.Entities()) == 1
		isAction := c.Callback() != nil
		var action string
		if isCmd {
			action = c.Text()
			idx := strings.Index(action, " ")
			if idx > 0 {
				action = action[:idx]
			}
		} else if isAction {
			action = strings.TrimSpace(strings.Split(c.Callback().Data, "|")[0])
		}
		bot.log.Errorw("error",
			"chat_id", c.Chat().ID,
			"chat_type", c.Chat().Type,
			"user_id", c.Sender().ID,
			"user_name", c.Sender().Username,
			"is_cmd", isCmd,
			"is_action", isAction,
			"action", action,
			"size", len(c.Text()),
			"err", err)
	}
}

func (bot *Bot) sendMessage(c tele.Context) ([]letta.Message, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		if err := c.Notify(tele.Typing); err != nil {
			bot.log.Warn("Failed to send initial typing notification", "error", err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.Notify(tele.Typing); err != nil {
					bot.log.Warn("Failed to send typing notification", "error", err)
				}
			}
		}
	}()

	return bot.thera.SendMessage(ctx, c.Sender().ID, c.Text())
}

func (bot *Bot) handleText(c tele.Context) error {
	messages, err := bot.sendMessage(c)
	if err != nil {
		bot.LogError(err, c)
		return c.Send("Sorry, I encountered an error processing your message.")
	}

	for _, msg := range messages {
		responseText := ""

		switch msg.GetMessageType() {
		case letta.MessageTypeAssistant:
			responseText = msg.GetContent()
		case letta.MessageTypeReasoning:
			responseText = fmt.Sprintf("* %s", msg.GetContent())
		case letta.MessageTypeHiddenReasoning:
			responseText = fmt.Sprintf("… %s", msg.GetContent())
		case letta.MessageTypeToolCall:
			responseText = fmt.Sprintf("→ %s", msg.GetContent())
		case letta.MessageTypeToolReturn:
			responseText = msg.GetContent()
			if responseText != "" {
				responseText = fmt.Sprintf("← %s", responseText)
			}
		default:
			bot.log.Warn("received message: ", msg)
			responseText = fmt.Sprintf("? %s:\n%s", msg.GetMessageType(), msg.GetContent())
		}

		if responseText != "" {
			if err := c.Send(responseText); err != nil {
				bot.LogError(err, c)
				return err
			}
		}
	}

	return nil
}

func (bot *Bot) handleStart(c tele.Context) error {
	chat, err := bot.api.ChatByID(c.Chat().ID)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = bot.thera.CreateChat(ctx, chat.ID, chat.FirstName, chat.Bio)
	if err != nil {
		bot.LogError(err, c)
		return c.Send("Sorry, I encountered an error starting new chat.")
	}

	return nil
}
