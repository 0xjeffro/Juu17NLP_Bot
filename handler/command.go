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

func checkPermission(username string) bool {
	db := orm.GetConn()
	var user orm.Users
	res := db.Where(&orm.Users{
		UserName: username,
	}).First(&user)
	if res.Error != nil {
		return false
	} else {
		return true
	}
}

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
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
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
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
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
	case "ping":
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db := orm.GetConn()

		var lastRelpyUpdatedByProducer orm.Replies
		var lastRelpyUpdatedByConsumer orm.Replies

		db.Order("created_at desc").First(&lastRelpyUpdatedByProducer)
		db.Order("updated_at desc").First(&lastRelpyUpdatedByConsumer)

		var kv orm.KV
		db.Where(&orm.KV{Key: "ProducerLastRun"}).First(&kv)
		producerLastRun := kv.Value

		var kv2 orm.KV
		// why kv2?
		// If the object’s primary key has been set, then condition query wouldn’t cover the value of primary key
		// but use it as a ‘and’ condition.
		// see: https://gorm.io/docs/query.html
		db.Where(&orm.KV{Key: "ConsumerLastRun"}).First(&kv2)
		consumerLastRun := kv2.Value

		msg.Text = fmt.Sprintf("🕔Producer:\n	run@ %s\n	update@ %s\n 🕔Consumer:\n	run@ %s\n	update@ %s",
			producerLastRun,
			lastRelpyUpdatedByProducer.CreatedAt.Format("2006-01-02 15:04:05"),
			consumerLastRun,
			lastRelpyUpdatedByConsumer.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
