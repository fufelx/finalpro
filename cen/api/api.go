package api

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	file   = &os.File{}
	Logger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
)

func Newcom(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDKey).(string)
	if !ok {
		http.Error(w, "Request ID not found in context", http.StatusInternalServerError)
		return
	}

	text := r.URL.Query().Get("text")

	if strings.Contains(text, "qwerty") || strings.Contains(text, "йцукен") || strings.Contains(text, "zxvbnm") {
		w.WriteHeader(403)
		go LoggerPrint(getRealIP(r), requestID, http.StatusForbidden)
		return
	}

	w.WriteHeader(200)
	go LoggerPrint(getRealIP(r), requestID, http.StatusOK)
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
func RequestIDMiddleware(next http.Handler) http.Handler {
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
