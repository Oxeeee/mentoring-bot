package tgbot

import (
	"fmt"
	"log"
	"time"

	"github.com/Oxeeee/klenov-bot/db"
	"github.com/Oxeeee/klenov-bot/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		handleStartCommand(bot, message)
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
	case "report":
		handleReportCommand(bot, message)
	case "broadcast":
		handleBroadcastCommand(bot, message)

	default:
		return
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

	reply := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я буду каждый день в 18:00 приходить к тебе, спрашивать твой фидбэк за день и ждать обратного сообщения :)")
	bot.Send(reply)
}

func handleReportCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var user domain.User
	if err := db.DB.Where("username = ?", message.From.UserName).First(&user).Error; err != nil {
		log.Printf("User %v dont registred", message.From.UserName)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ты не зарегистрирован в системе. Попроси ментора @y0na24 тебя зарегестрировать")
		bot.Send(reply)
		return
	}

	location, _ := time.LoadLocation("Europe/Moscow")
	now := time.Now().In(location)
	target := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, location)
	if now.Before(target) {
		reply := tgbotapi.NewMessage(message.Chat.ID, "Пока рановато для отчета, поработай еще")
		bot.Send(reply)
	}

	args := message.CommandArguments()
	if args == "" {
		log.Printf("Handle report: args is empty")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Используй команду так: /report отчет")
		bot.Send(reply)
		return
	}

	msg := domain.Message{Content: args, UserID: user.ID}
	if err := db.DB.Model(&domain.Message{}).Create(&msg).Error; err != nil {
		log.Printf("Error occured while add message to db: %v", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Произошла внутренняя ошибка, попробуй снова. Если не пройдет — обратись к поддержке /support")
		bot.Send(reply)
		return
	}

	resendReportToAdmins(message.From.UserName, args)

	reply := tgbotapi.NewMessage(message.Chat.ID, "Твой отчет записан, ментор тебе скоро ответит.")
	bot.Send(reply)
}

func resendReportToAdmins(sender string, message string) {
	msg := tgbotapi.NewMessage(-1002441023269, fmt.Sprintf("<b>Пользователь @%v отправил отчет:</b>\n%v", sender, message))
	msg.ParseMode = "HTML"
	_, err := tgbot.Send(msg)
	if err != nil {
		log.Printf("Error while sending daily notifications: %v", err)
	}
}
