package main

import (
	"context"
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
	"time"

	"github.com/joho/godotenv"
)

// change favicon notes
// change build notes
// starting with relative paths during dev
// adding an absolute path during production
// hot realoding with air
// adding ability to provide meta data for pages
// document how to add meta tags
// add ability to change page titles

var templates *template.Template

func main() {

	args := os.Args
	_ = godotenv.Load()
	port := os.Getenv("PORT")

	if len(args) > 1 && args[1] == "--build" {
		absolutePath := ""
		if len(args) > 2 {
			absolutePath = args[2]
		}
		if absolutePath == "" {
			fmt.Println("No absolute path provided, defaulting to localhost:8080")
			absolutePath = "http://localhost:8080"
		}
		if absolutePath[len(absolutePath)-1:] == "/" {
			panic("The absolute path should not end with a '/'")
		}
		serverDone := make(chan bool)
		shutdownComplete := make(chan bool)
		go func() {
			if err := runDevServer(serverDone, shutdownComplete, port); err != nil {
				panic(err)
			}
		}()
		build(absolutePath, port)
		serverDone <- true
		<-shutdownComplete
		runStaticServer(port)
		return
	}

	runDevServer(nil, nil, port)

}

func build(absolutePath string, port string) {
	cnf := config.GetDocConfig()
	filewriter.GenerateStaticAssets(cnf, absolutePath, port)
}

func runStaticServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		middleware.Chain(w, r, nil, func(cc *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
			handler.ServeOutFiles(w, r)
		})
	})
	fmt.Println("Static server is running on port: " + port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println(err)
	}
}

func runDevServer(serverDone chan bool, shutdownComplete chan bool, port string) error {
	mux := http.NewServeMux()
	parseTemplates()
	mux.HandleFunc("GET /favicon.ico", handler.ServeFavicon)
	mux.HandleFunc("GET /static/", handler.ServeStaticFiles)
	cnf := contentrouter.GenerateRoutes(mux, templates)
	filewriter.GenerateDynamicNavbar(cnf)
	server := &http.Server{Addr: ":" + port, Handler: mux}
	go func() {
		if serverDone != nil {
			<-serverDone
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			fmt.Println("Development server is shutting down...")
			if shutdownComplete != nil {
				shutdownComplete <- true // Signal that the shutdown is complete
			}
		}
	}()
	fmt.Println("Development server is running on port:", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func parseTemplates() {
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
}
