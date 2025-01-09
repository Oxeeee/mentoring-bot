package tgbot

import (
	"fmt"
	"log"

	"github.com/Oxeeee/klenov-bot/db"
	"github.com/Oxeeee/klenov-bot/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		go handleStartCommand(bot, message)
	case "adduser":
		go handleAddUserCommand(bot, message)
	case "removeuser":
		go handleRemoveUserCommand(bot, message)
	case "userlist":
		go handleUserListCommand(bot, message)
	case "addadmin":
		go handleAddAdminRightsCommand(bot, message)
	case "deleteadmin":
		go handleDeleteAdminRightsCommand(bot, message)

	default:
		reply := tgbotapi.NewMessage(message.Chat.ID, "Команда не распознана. Попробуй /start")
		bot.Send(reply)
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var user domain.User
	if err := db.DB.Where("username = ?", message.From.UserName).First(&user).Error; err != nil {
		log.Printf("User %v dont registred", message.From.UserName)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ты не зарегестрирован в системе. Попроси ментора @y0na24 тебя зарегестрировать")
		bot.Send(reply)
		return
	}

	if err := db.DB.Model(&domain.User{}).Where("username = ?", message.From.UserName).Update("chat_id", fmt.Sprintf("%v", message.Chat.ID)).Error; err != nil {
		log.Printf("Error occured while saving chat id: %v", err)
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я буду каждый день, в 18:00 приходить к тебе, спрашивать твой фидбэк за день, и ждать обратного сообщения :)")
	bot.Send(reply)
}
