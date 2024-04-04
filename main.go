package main

import (
	"godocument/internal/config"
	"godocument/internal/contentrouter"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", config.ServeFavicon)
	mux.HandleFunc("GET /static/", config.ServeStaticFiles)
	contentrouter.GenerateRoutes(mux)
	// err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	// if err != nil {
	// 	fmt.Println(err)
	// }

}
