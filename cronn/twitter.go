package cronn

import (
	"Juu17NLP_Bot/orm"
	"Juu17NLP_Bot/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Producer() {
	zap.S().Infof("Producer running...")

	scraper := twitterscraper.New().WithReplies(true)
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	envCookies := os.Getenv("COOKIE")
	var cookies []*http.Cookie
	err := json.NewDecoder(bytes.NewReader([]byte(envCookies))).Decode(&cookies)
	if err != nil {
		zap.S().Error(err.Error())
	}
	scraper.SetCookies(cookies)

	if scraper.IsLoggedIn() {
		zap.S().Infof("Twitter login successfully.")
	}
	db := orm.GetConn()
	for tweet := range scraper.SearchTweets(context.Background(), "(to:0xjuu_17)", 50) {
		if tweet.Error != nil {
			fmt.Println(err.Error())
		}

		replyID := tweet.ID
		replyID64, _ := strconv.ParseInt(replyID, 10, 64)
		author := tweet.Username
		text := tweet.Text
		url := tweet.PermanentURL

		var reply orm.Replies
		res := db.Where(orm.Replies{ReplyID: replyID64}).First(&reply)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				reply = orm.Replies{
					ReplyID: replyID64,
					Author:  author,
					Text:    text,
					Url:     url,
					Visited: false,
				}
				db.Create(&reply)
			}
		} else {
			zap.S().Infow("Tweet has existed.", "URL", url)
		}
	}

	// update status check
	var cstZone = time.FixedZone("CST", 8*3600)
	kv := orm.KV{
		Key:   "ProducerLastRun",
		Value: time.Now().In(cstZone).Format("2006-01-02 15:04:05"),
	}
	db.Save(&kv)
}

func Consumer(bot *tgbotapi.BotAPI) {
	var users []orm.Users
	var replies []orm.Replies
	zap.S().Infof("Consumer running...")
	db := orm.GetConn()
	db.Limit(20).Find(&users)
	db.Limit(30).Where("visited = ?", false).Find(&replies)

	for _, reply := range replies {
		if reply.Author == "0xjuu_17" {
			res := db.Model(&orm.Replies{}).Where(orm.Replies{ReplyID: reply.ReplyID}).
				Updates(orm.Replies{
					Visited:      true,
					PositiveProb: 0,
					NegativeProb: 0,
				})
			if res.Error != nil {
				zap.S().Error(res.Error.Error())
			}
			continue
		}

		// Check if the reply contains keywords
		keywords := utils.KeywordsAnalysis(reply.Text)

		// Check if the reply match regular expressions
		regex := utils.RegularExpressionAnalysis(reply.Text)

		// Sentiment analysis
		data := utils.SentimentAnalysis(reply.Text)
		positiveProb := data.Result.PositiveProb
		negativeProb := data.Result.NegativeProb

		// Concatenate message
		message := ""
		if len(keywords) > 0 {
			keywordsStr := strings.Join(keywords, ",")
			message += fmt.Sprintf("⚠️%s\n", keywordsStr)
		}

		if len(regex) > 0 {
			regexStr := strings.Join(regex, ",")
			message += fmt.Sprintf("📜%s\n", regexStr)
		}

		if negativeProb > 0.6 {
			message += fmt.Sprintf("🔥%s\n", "Negative Prob: "+strconv.FormatFloat(negativeProb, 'f', 3, 64))
		}

		if message != "" {
			message += fmt.Sprintf("%s\n%s\n", reply.Text, reply.Url)
			for _, user := range users {
				msg := tgbotapi.NewMessage(user.ChatID, message)
				_, err := bot.Send(msg)
				if err != nil {
					zap.S().Error(err.Error())
				}
			}
		}

		// Update
		res := db.Model(&orm.Replies{}).Where(orm.Replies{ReplyID: reply.ReplyID}).
			Updates(orm.Replies{
				Visited:      true,
				PositiveProb: positiveProb,
				NegativeProb: negativeProb,
			})
		if res.Error != nil {
			zap.S().Error(res.Error.Error())
		}
	}

	// update status check
	var cstZone = time.FixedZone("CST", 8*3600)
	kv := orm.KV{
		Key:   "ConsumerLastRun",
		Value: time.Now().In(cstZone).Format("2006-01-02 15:04:05"),
	}
	db.Save(&kv)
}
