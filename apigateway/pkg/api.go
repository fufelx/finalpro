package APIGateway

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func News(w http.ResponseWriter, r *http.Request) {
	// Основной сервис работает на `localhost` по умолчанию на порту 80 или 8080
	targetURL := "http://localhost:4040/news/40"

	// Отображаем информацию о запросе
	fmt.Printf("Проксирование запроса: %s %s\n", r.Method, targetURL)

	// Создаем новый запрос для целевого сервиса
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из оригинального запроса в проксируемый
	req.Header = r.Header

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса к целевому сервису", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ от целевого сервиса
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения ответа от сервиса", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ клиенту с оригинальным статусом
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func Filter(w http.ResponseWriter, r *http.Request) {
	reqid := ReqId()
	name := r.URL.Query().Get("name")
	pagest := r.URL.Query().Get("page")
	// Основной сервис работает на `localhost` по умолчанию на порту 80 или 8080
	targetURL := "http://localhost:4040/news?name=" + name + "&page=" + pagest + "&request_id=" + reqid

	// Отображаем информацию о запросе
	fmt.Printf("Проксирование запроса: %s %s\n", r.Method, targetURL)

	// Создаем новый запрос для целевого сервиса
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из оригинального запроса в проксируемый
	req.Header = r.Header

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса к целевому сервису", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ от целевого сервиса
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения ответа от сервиса", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ клиенту с оригинальным статусом
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func NewsFullDetailed(w http.ResponseWriter, r *http.Request) {
	reqid := ReqId()
	name := r.URL.Query().Get("newsid")
	var Renews, Recom []byte

	var wg sync.WaitGroup
	var mu sync.Mutex // Мьютекс для безопасного доступа к переменным Renews и Recom

	// Первая горутина для запроса комментариев
	wg.Add(1)
	go func() {
		defer wg.Done()

		targetURL := "http://localhost:4041/allcom?newsid=" + name + "&request_id=" + reqid
		fmt.Printf("Проксирование запроса: %s %s\n", r.Method, targetURL)

		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
			return
		}

		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Ошибка выполнения запроса к целевому сервису", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения ответа от сервиса", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		Recom = body
		mu.Unlock()
		log.Println("Получены данные Recom:", string(Recom))
	}()

	// Вторая горутина для запроса новостей
	wg.Add(1)
	go func() {
		defer wg.Done()

		targetURL := "http://localhost:4040/news?id=" + name + "&request_id=" + reqid
		fmt.Printf("Проксирование запроса: %s %s\n", r.Method, targetURL)

		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
			return
		}

		req.Header = r.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Ошибка выполнения запроса к целевому сервису", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения ответа от сервиса", http.StatusInternalServerError)
			return
		}

		mu.Lock()
		Renews = body
		mu.Unlock()
		log.Println("Получены данные Renews:", string(Renews))
	}()

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Проверяем, что оба ответа получены
	if Renews != nil && Recom != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(append(Renews, Recom...))
	} else {
		if Renews != nil {
			w.Write(Renews)
		} else {
			http.Error(w, "Ошибка: не удалось получить все необходимые данные", http.StatusInternalServerError)
		}
	}
}

func Comment(w http.ResponseWriter, r *http.Request) {
	reqid := ReqId()

	parentsID := r.URL.Query().Get("parentsid")
	newsID := r.URL.Query().Get("newsid")
	com := r.URL.Query().Get("com")

	targetURL := "http://localhost:4042/validate?text=" + strings.ReplaceAll(com, " ", "%20") + "&request_id=" + reqid
	fmt.Printf("Проксирование запроса: %s %s\n", r.Method, targetURL)

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса к целевому сервису", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения ответа от сервиса", http.StatusInternalServerError)
		return
	}
	log.Println(resp.Status)
	if resp.Status == "403 Forbidden" {
		http.Error(w, "Не прошел валидацию", http.StatusForbidden)
		return
	}

	targetURL = "http://localhost:4041/newcom?newsid=" + newsID + "&parentsid=" + parentsID + "&com=" + strings.ReplaceAll(com, " ", "%20") + "&request_id=" + reqid
	fmt.Printf("Проксирование запроса: %s %s\n", r.Method, targetURL)

	req, err = http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка выполнения запроса к целевому сервису", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

}

func ReqId() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел
	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
