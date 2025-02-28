package crypto_fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	dataFormatter "ticker-pulse-news-bot/internal/pkg/data_formatter"

	"github.com/joho/godotenv"
)

type NewsFetcher struct {
	baseURL string
}

func NewNewsFetcher() *NewsFetcher {
	return &NewsFetcher{
		baseURL: "https://newsdata.io/api/1/latest",
	}
}

func (nf *NewsFetcher) FetchLastNews() ([]dataFormatter.NewsMap, error) {
	envPath := os.Getenv("ENV_FILE")
	if envPath == "" {
		envPath = ".env"
	}

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("[NEWS_FETCHER]: Ошибка загрузки .env")
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("[NEWS_FETCHER]: Отсутствует NEWS_API_KEY")
	}

	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("q", "криптовалюта")
	params.Add("language", "ru")

	finalURL := fmt.Sprintf("%s?%s", nf.baseURL, params.Encode())

	resp, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("[NEWS_FETCHER]: Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[NEWS_FETCHER]: Ошибка при чтении ответа: %v", err)
	}

	var data map[string]any
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	formattedData, err := dataFormatter.ProcessLastNews(data)
	if err != nil {
		return nil, fmt.Errorf("[NEWS_FETCHER]: Ошибка при форматировании данных: %v", err)
	}

	return formattedData, nil
}
