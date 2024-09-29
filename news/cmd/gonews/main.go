package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/pkg/api"
	"main/pkg/rss"
	"main/pkg/storage"
	"net/http"
	"time"
)

// Каналы обработки новостей и ошибок
var (
	postcn = make(chan []storage.Post, 10)
	errcn  = make(chan error, 10)
)

// Структура конфига
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	// инициализация БД
	db, err := storage.New()
	if err != nil {
		errcn <- err
	}
	api := api.New(db)
	//Запуск рутин БД и ОШИБОК
	go AddBD(db)
	go ErrChecker()
	// чтение и раскодирование файла конфигурации
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		errcn <- err
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		errcn <- err
	}

	for _, url := range config.URLS {
		go parseURL(url, config.Period)
	}

	//Запуск localhost
	log.Println("Сервер запущен")
	err = http.ListenAndServe(":4040", api.Router())
	if err != nil {
		errcn <- err
	}
}

// rss + рутины добавления в БД и обработка ошибок
func parseURL(url string, period int) {
	for {
		posts, err := rss.Parse(url)
		if err != nil {
			errcn <- err
			continue
		}
		postcn <- posts
		time.Sleep(time.Minute * time.Duration(period))
	}
}

// функция добавления в бд
func AddBD(db *storage.DB) {
	for {
		select {
		case post := <-postcn:
			err := db.StoreNews(post)
			if err != nil {
				errcn <- err
			}
		}
	}
}

// функция обработки ошибок
func ErrChecker() {
	for {
		select {
		case err := <-errcn:
			log.Println(err)
		}
	}
}
