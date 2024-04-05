package main

import (
	"fmt"
	"godocument/internal/contentrouter"
	"godocument/internal/handler"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

var templates *template.Template

func main() {

	_ = godotenv.Load()

	var err error
	templates, err = template.ParseGlob("./html/templates/*.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", handler.ServeFavicon)
	mux.HandleFunc("GET /static/", handler.ServeStaticFiles)

	contentrouter.GenerateRoutes(mux, templates)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on port: " + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println(err)
	}

}
