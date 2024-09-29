package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"main/pkg/storage"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	file   = &os.File{}
	Logger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
)

type res struct {
	Pagi storage.Pagi   `json:"pagi"`
	News []storage.Post `json:"news"`
}
type API struct {
	db *storage.DB
	r  *mux.Router
}

// Конструктор API.
func New(db *storage.DB) *API {
	a := API{db: db, r: mux.NewRouter()}
	a.endpoints()
	return &a
}

// Router возвращает маршрутизатор для использования
// в качестве аргумента HTTP-сервера.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions)

	api.r.Use(requestIDMiddleware)

	api.r.HandleFunc("/news", api.findBybName).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)
	news, err := api.db.News(n)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}

func (api *API) findBybName(w http.ResponseWriter, r *http.Request) {

	requestID, ok := r.Context().Value(requestIDKey).(string)
	if !ok {
		http.Error(w, "Request ID not found in context", http.StatusInternalServerError)
		return
	}

	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")
	pagest := r.URL.Query().Get("page")

	if id != "" {
		idi, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			go LoggerPrint(getRealIP(r), requestID, http.StatusInternalServerError)

			return
		}
		news, err := api.db.NewsById(idi)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			go LoggerPrint(getRealIP(r), requestID, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
		response, err := json.Marshal(news)
		if err != nil {
			log.Fatal("ошибка кодирования json: ", err)
		}
		go LoggerPrint(getRealIP(r), requestID, http.StatusOK)
		w.Write(response)
		return
	}

	news, err := api.db.NewsByName(name)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		go LoggerPrint(getRealIP(r), requestID, http.StatusInternalServerError)
		return
	}

	maxPage := len(news) / 10
	if len(news)%10 != 0 {
		maxPage += 1
	}
	if len(news)/10 == 0 {
		maxPage = 1
	}

	resp := []res{}
	page := 1
	if pagest != "" {
		page, err = strconv.Atoi(pagest)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			go LoggerPrint(getRealIP(r), requestID, http.StatusInternalServerError)

			return
		}
	}

	if len(news) <= 10 {
		if page > maxPage {
			http.Error(w, "Крайняя страница - "+strconv.Itoa(maxPage), http.StatusInternalServerError)
			go LoggerPrint(getRealIP(r), requestID, http.StatusInternalServerError)
			return
		} else {
			pag := storage.Pagi{Pages: maxPage, CurrentPage: page, AmountOfElement: len(news)}
			resp = append(resp, res{Pagi: pag, News: news})
		}
	} else {
		if page > maxPage {
			http.Error(w, "Крайняя страница - "+strconv.Itoa(maxPage), http.StatusInternalServerError)
			go LoggerPrint(getRealIP(r), requestID, http.StatusInternalServerError)
			return
		} else {
			pag := storage.Pagi{Pages: maxPage, CurrentPage: page}
			if len(news) < page*10 {
				news = news[(page-1)*10:]
			} else {
				news = news[(page-1)*10 : 10*page]
			}
			pag = storage.Pagi{Pages: pag.Pages, CurrentPage: pag.CurrentPage, AmountOfElement: len(news)}
			resp = append(resp, res{Pagi: pag, News: news})
		}
	}
	w.WriteHeader(200)
	response, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("ошибка кодирования json: ", err)
	}
	go LoggerPrint(getRealIP(r), requestID, http.StatusOK)
	w.Write(response)
}

func LoggerPrint(ip, id string, code int) {
	file, _ = os.OpenFile("./logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	Logger = log.New(file, "", log.Ldate|log.Ltime)
	Logger.Println(ip, code, id)
}

func getRealIP(r *http.Request) string {
	// Попробуем получить IP из заголовка X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Если в X-Forwarded-For несколько IP-адресов, берём первый
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Если X-Forwarded-For отсутствует, попробуем заголовок X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// В противном случае, используем r.RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Если не удается разделить адрес и порт
	}

	return ip
}

type key string

const requestIDKey = key("request_id")

// Middleware для извлечения request_id из заголовков или параметров URL
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем request_id из параметров URL или заголовков
		requestID := r.URL.Query().Get("request_id")
		if requestID == "" {
			requestID = r.Header.Get("X-Request-ID")
		}

		// Сохраняем request_id в контексте
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		// Передаём новый контекст следующему обработчику
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
