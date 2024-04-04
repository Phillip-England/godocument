package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type CustomContext struct {
	context.Context
	StartTime time.Time
}

type CustomHandler func(cc *CustomContext, w http.ResponseWriter, r *http.Request)
type CustomMiddleware func(cc *CustomContext, w http.ResponseWriter, r *http.Request) error

func Chain(w http.ResponseWriter, r *http.Request, handler CustomHandler, middleware ...CustomMiddleware) {
	cc := &CustomContext{
		Context:   context.Background(),
		StartTime: time.Now(),
	}
	for _, mw := range middleware {
		err := mw(cc, w, r)
		if err != nil {
			return
		}
	}
	handler(cc, w, r)
	Log(cc, w, r)
}

func Log(cc *CustomContext, w http.ResponseWriter, r *http.Request) error {
	elapsedTime := time.Since(cc.StartTime)
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
	return nil
}
