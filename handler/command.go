package handler

import (
	"Juu17NLP_Bot/orm"
	"Juu17NLP_Bot/utils"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
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
				msg.Text = fmt.Sprintf("关键词已添加")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			msg.Text = fmt.Sprintf("关键词已存在")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	case "ar":
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
			Type:    "regex",
		}).First(&rule)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				rule = orm.Rules{
					Content: msgText,
					Type:    "regex",
				}
				db.Create(&rule)
				msg.Text = fmt.Sprintf("正则已添加")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		} else {
			msg.Text = fmt.Sprintf("正则已存在")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	case "delete_keyword":
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		msgText := update.Message.CommandArguments()
		targeId, err := strconv.Atoi(msgText)
		if err != nil {
			msg.Text = fmt.Sprintf("Error: 请输入要删除的关键词id")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db := orm.GetConn()
		var rule []orm.Rules
		db.Limit(300).Where(&orm.Rules{Type: "keyword"}).Find(&rule)
		if targeId > len(rule) {
			msg.Text = fmt.Sprintf("Error: 请输入正确的id")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db.Delete(&rule[targeId-1])
		msg.Text = fmt.Sprintf("关键词已删除")
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "delete_regex":
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		msgText := update.Message.CommandArguments()
		targeId, err := strconv.Atoi(msgText)
		if err != nil {
			msg.Text = fmt.Sprintf("Error: 请输入要删除的正则id")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db := orm.GetConn()
		var rule []orm.Rules
		db.Limit(300).Where(&orm.Rules{Type: "regex"}).Find(&rule)
		if targeId > len(rule) {
			msg.Text = fmt.Sprintf("Error: 请输入正确的id")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db.Delete(&rule[targeId-1])
		msg.Text = fmt.Sprintf("正则已删除")
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "list_keywords":
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db := orm.GetConn()
		var rules []orm.Rules
		db.Limit(300).Where(&orm.Rules{Type: "keyword"}).Find(&rules)
		keywords := ""
		for idx, rule := range rules {
			keywords += fmt.Sprintf("%d. %s\n", idx+1, rule.Content)
		}
		msg.Text = fmt.Sprintf("所有关键词：\n" + keywords)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "list_regex":
		if checkPermission(update.Message.From.UserName) == false {
			msg.Text = fmt.Sprintf("无权限")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			return
		}
		db := orm.GetConn()
		var rules []orm.Rules
		db.Limit(300).Where(&orm.Rules{Type: "regex"}).Find(&rules)
		regex := ""
		for idx, rule := range rules {
			regex += fmt.Sprintf("%d. %s\n", idx+1, rule.Content)
		}
		msg.Text = fmt.Sprintf("所有正则：\n" + regex)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case "test":
		msgText := update.Message.CommandArguments()
		keyword := utils.KeywordsAnalysis(msgText)
		regex := utils.RegularExpressionAnalysis(msgText)
		msg.Text = fmt.Sprintf("关键词：%s\n正则：%s", keyword, regex)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
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
