package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func InitBot(token string, webhook string, debug bool) *tgbotapi.BotAPI {
	var err error
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	wh, _ := tgbotapi.NewWebhook(webhook)
	wh.AllowedUpdates = []string{"chat_member", "my_chat_member", "callback_query",
		"channel_post", "edited_channel_post", "edited_message", "inline_query", "message",
		"pre_checkout_query", "shipping_query"}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	return bot
}

func InitLogger() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig = encoderConfig
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

// CheckEnv 检查必要的环境变量是否设置
func CheckEnv() {
	if os.Getenv("BOT_TOKEN") == "" {
		panic("BOT_TOKEN is not set")
	}
	if os.Getenv("WEBHOOK") == "" {
		panic("WEBHOOK is not set")
	}
	if os.Getenv("PORT") == "" {
		panic("PORT is not set")
	}
	// 该bot只服务于一个群组，因此这里需要在环境变量中配置当前服务的群组ID
	if os.Getenv("SPELL") == "" {
		panic("SPELL is not set")
	}
}
