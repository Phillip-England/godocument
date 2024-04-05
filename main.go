package main

import (
	"fmt"
	"godocument/internal/contentrouter"
	"godocument/internal/handler"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", handler.ServeFavicon)
	mux.HandleFunc("GET /static/", handler.ServeStaticFiles)
	contentrouter.GenerateRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := http.ListenAndServe(":"+os.Getenv(port), mux)
	if err != nil {
		fmt.Println(err)
	}

}
