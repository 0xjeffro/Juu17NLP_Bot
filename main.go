package main

import (
	"Juu17NLP_Bot/cronn"
	"Juu17NLP_Bot/handler"
	"Juu17NLP_Bot/orm"
	"Juu17NLP_Bot/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var bot *tgbotapi.BotAPI

func main() {
	utils.CheckEnv()
	orm.AutoCreateTable()
	token := func() string {
		if os.Getenv("BOT_TOKEN") == "" {
			panic("BOT_TOKEN is not set")
		} else {
			return os.Getenv("BOT_TOKEN")
		}
	}()
	webhook := func() string {
		if os.Getenv("WEBHOOK") == "" {
			panic("WEBHOOK is not set")
		}
		return strings.TrimSuffix(os.Getenv("WEBHOOK"), "/")
	}()
	port := func() string {
		if os.Getenv("PORT") == "" {
			return "8080"
		}
		return os.Getenv("PORT")
	}()
	debug := os.Getenv("DEBUG") == "true"
	utils.InitLogger()

	webhookSuffix := utils.MD5(token)
	bot = utils.InitBot(token, webhook+"/"+webhookSuffix, debug)

	producer := cron.New(cron.WithSeconds())
	_, err := producer.AddFunc("@every 360s", func() {
		cronn.Producer()
	})
	if err != nil {
		log.Println(err)
	}
	producer.Start()

	consumer := cron.New(cron.WithSeconds())
	_, err = consumer.AddFunc("@every 180s", func() {
		cronn.Consumer(bot)
	})
	if err != nil {
		log.Println(err)
	}
	consumer.Start()

	startGin(webhookSuffix, port, debug)

}

func startGin(webhookSuffix string, port string, debug bool) {

	router := gin.New()
	router.Use(utils.Cors())
	if debug {
		router.Use(gin.Logger())
	}
	router.POST("/"+webhookSuffix, webhookHandler)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}

func webhookHandler(c *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(c.Request.Body)

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.Message != nil {
		zap.S().Infow("Received a message.",
			"chat_id", update.Message.Chat.ID,
			"message_id", update.Message.MessageID,
			"from", update.Message.From,
			"first_name", update.Message.From.FirstName,
			"last_name", update.Message.From.LastName,
			"text", update.Message.Text,
			"date", update.Message.Date,
		)
		if update.Message.IsCommand() {
			handler.CommandHandler(bot, update)
		}
	}
}
