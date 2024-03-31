package main

import (
	"fmt"
	"godocument/internal/config"
	"godocument/internal/contentrouter"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()

	contentrouter.GenerateRoutes()

	mux.HandleFunc("GET /favicon.ico", config.ServeFavicon)
	mux.HandleFunc("GET /static/", config.ServeStaticFiles)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}
