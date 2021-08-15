package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/MrSedan/gotgbot/dispatcher"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var Hub = &dispatcher.Hub{
	Servers:   make(map[int64]*dispatcher.DispServer),
	NewServer: make(chan *dispatcher.DispServer, 100),
	RemServer: make(chan *dispatcher.DispServer, 100),
}

func init() {
	// Загрузка ключа из .env файла
	if err := godotenv.Load(); err != nil {
		log.Print("No env file")
	}
	go Hub.Run()
}

func main() {
	var cond = true
	key, exists := os.LookupEnv("API_KEY")
	if !exists {
		log.Panic("No api key!")
	}
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		log.Panic(err)
	}
	time.Sleep(time.Millisecond * 500)
	for len(updates) != 0 {
		<-updates
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if strings.ToLower(update.Message.Text) == "карась" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Starting...")
			bot.Send(msg)
			cond = true
			continue
		}
		if !cond {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if strings.ToLower(update.Message.Text) == "щука" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Stopping...")
			bot.Send(msg)
			cond = false
			continue
		}

		s, ok := Hub.Servers[update.Message.Chat.ID]
		if !ok {
			s = dispatcher.NewServer(update.Message.Chat.ID, Hub)
			Hub.NewServer <- s
			go s.Run()
		}
		d, ok := s.Disps[update.Message.From.ID]
		if !ok {
			d = dispatcher.CreateDisp(update.Message.From.ID, s, dispatcher.HelloHandler())
			s.NewDisp <- d
			go d.Run(bot)
		}
		d.NextStage <- true
	}
}
