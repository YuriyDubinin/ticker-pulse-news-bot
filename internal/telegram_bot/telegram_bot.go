package telegram_bot

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	godotenv "github.com/joho/godotenv"
)

type TelegramBot struct {
	api    *tgbotapi.BotAPI
	chatID string
}

func NewTelegramBot() (*TelegramBot, error) {
	envPath := os.Getenv("ENV_FILE")
	if envPath == "" {
		envPath = ".env" // значение по умолчанию
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("[TICKER-PULSE-NEWS-BOT]: Ошибка загрузки .env")
	}

	apiKey := os.Getenv("TELEGRAM_BOT_API_KEY")
	chatID := os.Getenv("TELEGRAM_GROUP_ID")

	if apiKey == "" || chatID == "" {
		log.Fatal("[TICKER-PULSE-NEWS-BOT]: TELEGRAM_BOT_API_KEY / TELEGRAM_GROUP_ID не найдено в .env")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}

	log.Printf("[TICKER-PULSE-NEWS-BOT]: Авторизован как %s\n", bot.Self.UserName)

	return &TelegramBot{
		api:    bot,
		chatID: chatID,
	}, nil
}

func (tb *TelegramBot) SendMessageToChannel(text string) error {
	msg := tgbotapi.NewMessageToChannel(tb.chatID, text)
	_, err := tb.api.Send(msg)
	return err
}
