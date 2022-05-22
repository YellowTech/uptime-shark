package notifications

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateTelegramService(apiKey string, chatId int64, name string) (*NotificationService, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	
    if err != nil {
		return nil, err
    }

	telegramBot := new(NotificationService)
	telegramBot.Name = name
	telegramBot.SendMessage = func(message string) error {return sendMessage(bot, chatId, message)}

	return telegramBot, err
}

func sendMessage(bot *tgbotapi.BotAPI, chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := bot.Send(msg)
	return err
}
