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

func handleAddUserCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := message.CommandArguments()
	if args == "" {
		log.Printf("Add user: args is empty")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Используй команду так: /adduser username")
		bot.Send(reply)
		return
	}

	var admin domain.User
	db.DB.Where("username = ?", message.From.UserName).First(&admin)
	if admin.Role != "admin" {
		log.Printf("User has not admin rights: %v", admin.Username)
		reply := tgbotapi.NewMessage(message.Chat.ID, "У тебя нет прав для выполнения этой команды.")
		bot.Send(reply)
		return
	}

	newUser := domain.User{Username: args, IsWhitelisted: true, Role: "user"}
	if err := db.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error while creating new user: %v", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при добавлении пользователя.")
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Пользователь "+args+" добавлен.")
	bot.Send(reply)
}

func handleRemoveUserCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	args := message.CommandArguments()
	if args == "" {
		log.Printf("Remove user: args is empty")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Используй команду так: /adduser username")
		bot.Send(reply)
		return
	}

	var admin domain.User
	db.DB.Where("username = ?", message.From.UserName).First(&admin)
	if admin.Role != "admin" {
		log.Printf("User has not admin rights: %v", admin.Username)
		reply := tgbotapi.NewMessage(message.Chat.ID, "У тебя нет прав для выполнения этой команды.")
		bot.Send(reply)
		return
	}

	if err := db.DB.Where("username = ?", args).Delete(&domain.User{}).Error; err != nil {
		log.Printf("Error while deleting user: %v", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при удалении пользователя.")
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Пользователь "+args+" удален.")
	bot.Send(reply)
}
