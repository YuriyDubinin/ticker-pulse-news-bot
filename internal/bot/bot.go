package bot

import (
	"fmt"
	"log"
	newsFetcher "ticker-pulse-news-bot/internal/news_fetcher"
	dataFormatter "ticker-pulse-news-bot/internal/pkg/data_formatter"
	telegramBot "ticker-pulse-news-bot/internal/telegram_bot"
	workerPool "ticker-pulse-news-bot/internal/worker_pool"
	"time"
)

type Bot struct {
	tgBot      *telegramBot.TelegramBot
	workerPool *workerPool.WorkerPool
	news       []dataFormatter.NewsMap
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
func (b *Bot) Start() {
	log.Println("[TICKER-PULSE-NEWS-BOT]: Бот запущен")
	b.workerPool.Start()
	b.tgBot.SendMessageToChannel("🌐 Новости меняются каждую секунду")
	b.UpdateNews()
	b.PostNews()
}

// Остановка WorkerPool
func (b *Bot) Stop() {
	b.workerPool.Stop()
	log.Println("[TICKER-PULSE-NEWS-BOT]: Бот остановлен")
}

// Временное решение до подъема базы
func (b *Bot) UpdateNews() {
	nf := newsFetcher.NewNewsFetcher()
	b.workerPool.AddTask(func() {
		for {
			news, err := nf.FetchLastNews()
			if err != nil {
				log.Fatalf("[TICKER-PULSE-NEWS-BOT]: Ошибка при получении новостей: %v", err)
			}

			b.news = news
			log.Printf("[TICKER-PULSE-NEWS-BOT]: Новости успешно обновлены")
			time.Sleep(time.Duration(5 * time.Hour))
		}
	})
}

func (b *Bot) PostNews() {
	time.Sleep(time.Duration(2 * time.Minute))
	b.workerPool.AddTask(func() {
		for {
			for _, item := range b.news {
				b.tgBot.SendMessageToChannel(fmt.Sprintf("%v. %v %v", item.Title, item.Description, item.Link))
				time.Sleep(time.Duration(30 * time.Minute))
			}
		}
	})
}
