package actions

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// DeleteMessage 删除消息
func DeleteMessage(bot *tgbotapi.BotAPI, chatId int64, messageId int) {
	_, err := bot.Request(tgbotapi.DeleteMessageConfig{
		ChatID:    chatId,
		MessageID: messageId,
	})
	if err != nil {
		return
	}
}

func SendTextMessage(bot *tgbotapi.BotAPI, chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}

// PinMessage 消息置顶
func PinMessage(bot *tgbotapi.BotAPI, chatId int64, messageId int, notification bool) {
	_, err := bot.Request(tgbotapi.PinChatMessageConfig{
		ChatID:              chatId,
		MessageID:           messageId,
		DisableNotification: !notification,
	})
	if err != nil {
		return
	}
}

// SendInlineButtonMessage 发送带有inline button的消息
func SendInlineButtonMessage(bot *tgbotapi.BotAPI, chatId int64, text string, inlineKeyboard [][]tgbotapi.InlineKeyboardButton) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard...)
	_, err := bot.Send(msg)
	if err != nil {
		return
	}
}
