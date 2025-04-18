package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Hello World")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println(err)
		return
	}

	var chatIdString = os.Getenv("USER_ID")
	chatId, err := strconv.ParseInt(chatIdString, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := tgbotapi.NewMessage(chatId, "Hello World")
	bot.Send(msg)
}
