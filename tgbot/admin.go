package tgbot

import (
	"fmt"
	"log"

	"github.com/Oxeeee/klenov-bot/db"
	"github.com/Oxeeee/klenov-bot/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkAdmin(bot *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	var admin domain.User
	db.DB.Where("username = ?", message.From.UserName).First(&admin)
	if admin.Role != "admin" {
		log.Printf("User has not admin rights: %v", admin.Username)
		reply := tgbotapi.NewMessage(message.Chat.ID, "У тебя нет прав для выполнения этой команды.")
		bot.Send(reply)
		return false
	}
	return true
}

func handleAddUserCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if admin := checkAdmin(bot, message); admin == false {
		return
	}

	args := message.CommandArguments()
	if args == "" {
		log.Printf("Add user: args is empty")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Используй команду так: /adduser username")
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

	if admin := checkAdmin(bot, message); admin == false {
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

func handleUserListCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if admin := checkAdmin(bot, message); admin == false {
		return
	}

	var users []string
	if err := db.DB.Model(&domain.User{}).Pluck("username", &users).Error; err != nil {
		log.Printf("Error while listing users: %v", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при поиске пользователей.")
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Пользователи бота: %v", users))
	bot.Send(reply)
}

func handleAddAdminRightsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if admin := checkAdmin(bot, message); admin == false {
		return
	}

	args := message.CommandArguments()
	if args == "" {
		log.Printf("Add admin: args is empty")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Используй команду так: /addadmin username")
		bot.Send(reply)
		return
	}

	if err := db.DB.Model(&domain.User{}).Where("username = ?", args).Update("role", "admin").Error; err != nil {
		log.Printf("Error while adding admin rights: %v", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при добавлении администратора.")
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Пользователю "+args+" выданы права администратора.")
	bot.Send(reply)
}

func handleDeleteAdminRightsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if admin := checkAdmin(bot, message); admin == false {
		return
	}

	args := message.CommandArguments()
	if args == "" {
		log.Printf("Delete admin: args is empty")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Используй команду так: /deleteadmin username")
		bot.Send(reply)
		return
	}

	if err := db.DB.Model(&domain.User{}).Where("username = ?", args).Update("role", "user").Error; err != nil {
		log.Printf("Error while adding admin rights: %v", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ошибка при удалении администратора.")
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Пользователь "+args+" исключен из списка администраторов.")
	bot.Send(reply)
}
