package main

import (
	"fmt"
	"godocument/internal/contentrouter"
	"godocument/internal/filewriter"
	"godocument/internal/handler"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/joho/godotenv"
)

var templates *template.Template

func main() {

	args := os.Args

	_ = godotenv.Load()

	templates = template.New("")
	err := filepath.Walk("./html", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err := templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", handler.ServeFavicon)
	mux.HandleFunc("GET /static/", handler.ServeStaticFiles)

	cnf := contentrouter.GenerateRoutes(mux, templates)
	filewriter.GenerateDynamicNavbar(cnf)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if len(args) > 1 && args[1] == "--build" {
		filewriter.GenerateStaticAssets(cnf)
		return
	}

	fmt.Println("Server is running on port: " + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println(err)
	}

}
