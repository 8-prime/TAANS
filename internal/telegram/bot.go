package telegram

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message struct {
	Text string
}

type Config struct {
	BotToken string
	UserId   int64
}

func LoadConfig() (Config, error) {
	var config Config

	token, found := os.LookupEnv("BOT_TOKEN")
	if !found {
		fmt.Println(found)
		return config, fmt.Errorf("BOT_TOKEN not found")
	}
	config.BotToken = token

	userId, found := os.LookupEnv("USER_ID")
	if !found {
		fmt.Println(found)
		return config, fmt.Errorf("USER_ID not found")
	}
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return config, fmt.Errorf("USER_ID is not a valid integer")
	}
	config.UserId = userIdInt

	return config, nil
}

func StartBot(config Config, update chan Message) {
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		msg := <-update
		message := tgbotapi.NewMessage(config.UserId, msg.Text)
		_, err := bot.Send(message)
		if err != nil {
			fmt.Println(err)
		}
	}

}
