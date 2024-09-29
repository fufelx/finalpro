package main

import (
	APIGateway "APIGateway/pkg"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/news", APIGateway.News)
	http.HandleFunc("/filter", APIGateway.Filter)
	http.HandleFunc("/newsfulldetailed", APIGateway.NewsFullDetailed)
	http.HandleFunc("/comment", APIGateway.Comment)
	fmt.Println("GETWAY RUNING LOCALHOST")
	log.Fatal(http.ListenAndServe(":80", nil))
}
