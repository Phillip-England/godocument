package main

import (
	"fmt"
	"godocument/internal/config"
	"godocument/internal/contentrouter"
	"godocument/internal/filewriter"
	"godocument/internal/handler"
	"godocument/internal/middleware"
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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", handler.ServeFavicon)

	if len(args) > 1 && args[1] == "--reset" {
		filewriter.ResetOutDir()
		filewriter.ResetDocsDir()
		filewriter.ResetGodocumentConfig()
		return
	}

	if len(args) > 1 && args[1] == "--build" {
		port := os.Getenv("STATIC_PORT")
		if port == "" {
			port = "8000"
		}
		cnf := config.GetDocConfig()
		filewriter.GenerateStaticAssets(cnf)
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			middleware.Chain(w, r, nil, func(cc *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
				handler.ServeOutFiles(w, r)
			})
		})
		fmt.Println("Server is running on port: " + port)
		err := http.ListenAndServe(":"+port, mux)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

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

	mux.HandleFunc("GET /static/", handler.ServeStaticFiles)

	cnf := contentrouter.GenerateRoutes(mux, templates)
	filewriter.GenerateDynamicNavbar(cnf)

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
