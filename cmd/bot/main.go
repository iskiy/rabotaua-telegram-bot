package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iskiy/rabotaua-telegram-bot/internal/database"
	"github.com/iskiy/rabotaua-telegram-bot/internal/telegram"
	"github.com/iskiy/rabotaua-telegram-bot/pkg/rabotaua"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func getEnvField(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file\n")
	}
	return os.Getenv(key)
}

func main() {
	botAPI, err := tgbotapi.NewBotAPI(getEnvField("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	//botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)
	client := rabotaua.NewRabotaClient()
	storage, err := database.NewStorage(getEnvField("DB_PATH"))
	if err != nil {
		log.Fatal(err.Error())
	}
	schedulesMap, err := client.GetSchedulesMap()
	if err != nil {
		log.Fatal(err.Error())
	}
	bot := telegram.NewRabotaUABot(botAPI, client, storage, schedulesMap)
	err = bot.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
