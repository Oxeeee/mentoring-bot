package tgbot

import (
	"log"

	"github.com/Oxeeee/klenov-bot/db"
	"github.com/Oxeeee/klenov-bot/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		handleStartCommand(bot, message)
	case "whitelist":
		handleWhitelistCommand(bot, message)
	case "adduser":
		handleAddUserCommand(bot, message)
	case "removeuser":
		handleRemoveUserCommand(bot, message)
	case "userlist":
		handleUserListCommand(bot, message)
	case "addadmin":
		handleAddAdminRightsCommand(bot, message)
	case "deleteadmin":
		handleDeleteAdminRightsCommand(bot, message)

	default:
		reply := tgbotapi.NewMessage(message.Chat.ID, "Команда не распознана. Попробуй /start")
		bot.Send(reply)
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я буду каждый день, в 18:00 приходить к тебе, спрашивать твой фидбэк за день, и ждать обратного сообщения :)")
	bot.Send(reply)
}

func handleWhitelistCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var user domain.User
	res := db.DB.Where("username = ?", message.From.UserName).First(&user)

	if res.Error != nil {
		log.Printf("User %v dont registred", message.From.UserName)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ты не зарегестрирован в системе.")
		bot.Send(reply)
		return
	}

	if user.IsWhitelisted {
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ты находишься в белом списке.")
		bot.Send(reply)
	} else {
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ты не находишься в белом списке.")
		bot.Send(reply)
	}
}
