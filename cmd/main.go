package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	bot "ticker-pulse-news-bot/internal/bot"

	"github.com/joho/godotenv"
)

func main() {
	envPath := os.Getenv("ENV_FILE")
	if envPath == "" {
		envPath = ".env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("[TICKER-PULSE-NEWS-BOT]: Ошибка загрузки .env")
	}

	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		log.Fatal("[TICKER-PULSE-NEWS-BOT]: MAX_WORKERS должно быть числом")
	}

	tickerPulseBot, err := bot.NewBot(maxWorkers)
	if err != nil {
		log.Fatalf("[TICKER-PULSE-NEWS-BOT]: Ошибка инициализации бота: %v", err)
	}

	tickerPulseBot.Start()

	// Ожидание SIGINT/SIGTERM для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[TICKER-PULSE-NEWS-BOT]: Завершение работы..")
	tickerPulseBot.Stop()
}
