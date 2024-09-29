package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"main/api"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	// Подключаем middleware ко всему роутеру
	r.Use(api.RequestIDMiddleware)
	r.HandleFunc("/validate", api.Newcom).Methods(http.MethodGet, http.MethodOptions)

	fmt.Println("Цензура порт 4042")
	log.Fatal(http.ListenAndServe(":4042", r))
}
