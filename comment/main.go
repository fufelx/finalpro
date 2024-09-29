package main

import (
	"comment/pkg/api"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	if api.Errdb != nil {
		log.Fatal("ошибка подключения к БД: ", api.Errdb)
	}
	r := mux.NewRouter()

	// Подключаем middleware ко всему роутеру
	r.Use(api.RequestIDMiddleware)
	r.HandleFunc("/newcom", api.Newcom).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/allcom", api.Allcom).Methods(http.MethodGet, http.MethodOptions)

	fmt.Println("Server is running on port 4041")
	log.Fatal(http.ListenAndServe(":4041", r))
}
