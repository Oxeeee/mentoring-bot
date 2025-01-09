package tgbot

import (
	"log"
	"time"

	"github.com/Oxeeee/klenov-bot/db"
	"github.com/Oxeeee/klenov-bot/domain"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tgbot *tgbotapi.BotAPI

func sendNotification(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := tgbot.Send(msg)
	if err != nil {
		log.Printf("Error while sending daily notifications: %v", err)
	}
}

func sendDailyNotifications() {
	log.Println("Отправка ежедневных уведомлений...")

	var users []domain.User
	if err := db.DB.Model(&domain.User{}).Find(&users).Error; err != nil {
		log.Printf("Error while finding users: %v", err)
		return
	}

	for _, user := range users {
		sendNotification(user.ChatID, "Привет! Жду твой фидбэк за день!")
	}
}

func InitScheduler(bot *tgbotapi.BotAPI) {
	tgbot = bot

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("Error while loading timezone: %v", err)
	}
	// log.Println(time.Now())

	scheduler := gocron.NewScheduler(location)

	scheduler.Every(1).Day().At("18:00").Do(sendDailyNotifications)

	scheduler.StartAsync()
}
