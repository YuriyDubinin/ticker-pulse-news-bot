package data_formatter

import (
	"fmt"
	"log"
)

type NewsMap struct {
	ID          string
	Title       string
	Description string
	Link        string
	ImgURL      string
	VideoURL    string
}

func ProcessLastNews(lastNews map[string]any) ([]NewsMap, error) {
	var newsList []NewsMap

	results, ok := lastNews["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("[ProcessLastNews]: не удалось привести results к []interface{}")
	}

	for _, news := range results {
		newsMap, ok := news.(map[string]any)
		if !ok {
			log.Println("[ProcessLastNews]: Пропущен элемент: не является map[string]any")
			continue
		}

		ID, _ := newsMap["article_id"].(string)
		title, _ := newsMap["title"].(string)
		description, _ := newsMap["description"].(string)
		link, _ := newsMap["link"].(string)
		imgURL, _ := newsMap["img_url"].(string)
		videoURL, _ := newsMap["video_url"].(string)

		newsList = append(newsList, NewsMap{
			ID:          ID,
			Title:       title,
			Description: description,
			Link:        link,
			ImgURL:      imgURL,
			VideoURL:    videoURL,
		})
	}

	if len(newsList) == 0 {
		return nil, fmt.Errorf("[ProcessLastNews]: список новостей пуст")
	}

	return newsList, nil
}
