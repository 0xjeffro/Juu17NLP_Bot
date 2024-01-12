package handler

import (
	"Juu17NLP_Bot/orm"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func CommandHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cmd := update.Message.Command()
	zap.S().Info("Receive Command: \\" + cmd + ".")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch cmd {
	case "start":
		msg.Text = "👋🏻 Developed by Jeffro."
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true

		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "spell":
		// 获取消息内容
		msgText := update.Message.CommandArguments()
		// 获取消息发送者的ID
		// msgFrom := update.Message.From.ID
		Spell := os.Getenv("SPELL")
		if msgText == Spell {
			// 回复消息
			db := orm.GetConn()
			user := orm.Users{UserName: update.Message.From.UserName,
				ChatID: update.Message.From.ID}
			fmt.Println(time.Time{})
			db.Create(&user)
			msg.Text = fmt.Sprintf("身份验证通过")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		} else {
			msg.Text = fmt.Sprintf("验证失败")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}

	case "a":
		msgText := update.Message.CommandArguments()

		var rule orm.Rules

		db := orm.GetConn()
		res := db.Where(&orm.Rules{
			Content: msgText,
			Type:    "keyword",
		}).First(&rule)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				rule = orm.Rules{
					Content: msgText,
					Type:    "keyword",
				}
				db.Create(&rule)
				msg.Text = fmt.Sprintf("规则已添加")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			msg.Text = fmt.Sprintf("规则已存在")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	case "d":
		msgText := update.Message.CommandArguments()
		var rule orm.Rules

		db := orm.GetConn()
		res := db.Where(&orm.Rules{
			Content: msgText,
			Type:    "keyword",
		}).First(&rule)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				msg.Text = fmt.Sprintf("该规则不存在")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			db.Delete(&rule)
			msg.Text = fmt.Sprintf("规则已删除")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
