package bot

import (
	"log"
	telegramBot "ticker-pulse-news-bot/internal/telegram_bot"
	workerPool "ticker-pulse-news-bot/internal/worker_pool"
)

type Bot struct {
	tgBot      *telegramBot.TelegramBot
	workerPool *workerPool.WorkerPool
}

func NewBot(maxWorkers int) (*Bot, error) {
	tgBot, err := telegramBot.NewTelegramBot()
	if err != nil {
		log.Fatal("[TICKER-PULSE-NEWS-BOT]: Ошибка инициализации Telegram бота: ", err)
		return nil, err
	}

	wp := workerPool.NewWorkerPool(maxWorkers)

	return &Bot{
		tgBot:      tgBot,
		workerPool: wp,
	}, nil
}

// Запуск вместе с WorkerPool
// todo: на следующей итерации обеспечить гибкость currency
func (b *Bot) Start() {
	log.Println("[TICKER-PULSE-NEWS-BOT]: Бот запущен")
	b.workerPool.Start()
	b.SendMessageAsync("🌐 Новости меняются каждую секунду")

}

// Остановка WorkerPool
func (b *Bot) Stop() {
	b.workerPool.Stop()
	log.Println("[TICKER-PULSE-NEWS-BOT]: Бот остановлен")
}

// Отправка сообщений асинхронно через WorkerPool
func (b *Bot) SendMessageAsync(text string) {
	b.workerPool.AddTask(func() {
		if err := b.tgBot.SendMessageToChannel(text); err != nil {
			log.Printf("[TICKER-PULSE-NEWS-BOT]: Ошибка отправки сообщения: %v", err)
		} else {
			log.Println("[TICKER-PULSE-NEWS-BOT]: Сообщение успешно отправлено")
		}
	})
}

func (b *Bot) CheckQuoteLimitsByInterval(interval int) {
}
