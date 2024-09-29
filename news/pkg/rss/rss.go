package rss

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"main/pkg/storage"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
)

func Parse(rssURL string) ([]storage.Post, error) {
	// Выполняем HTTP-запрос для получения RSS-канала
	resp, err := http.Get(rssURL)
	if err != nil {
		log.Fatalf("Ошибка при запросе RSS-канала: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Ошибка: статус код %d %s", resp.StatusCode, resp.Status)
	}

	// Парсинг RSS-канала с помощью gofeed
	fp := gofeed.NewParser()
	rssFeed, err := fp.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при парсинге RSS: %v", err)
	}

	// Вывод информации о новостях
	var Posts []storage.Post
	for _, item := range rssFeed.Items {
		var i int64
		item.Published = strings.ReplaceAll(item.Published, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.Published)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.Published)
		}
		if err == nil {
			i = t.Unix()
		} else {
			i, err = strconv.ParseInt(item.Published, 10, 64)
			if err != nil {
				log.Fatalf("Ошибка при парсинге RSS: %v", err)
			}
		}
		Posts = append(Posts, storage.Post{Title: item.Title, Content: strip.StripTags(item.Description), PubTime: i, Link: item.Link})
	}
	return Posts, nil
}
