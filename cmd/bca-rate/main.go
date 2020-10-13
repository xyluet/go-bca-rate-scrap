package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-scrap/internal/scrap"
	"go-scrap/internal/scrap/colly"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	var collyScrapper scrap.Scrapper
	{
		collyScrapper = colly.NewScrapper(colly.DefaultFactory())
	}

	var service scrap.Service
	service = scrap.NewService(
		collyScrapper,
	)

	router := chi.NewRouter()
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rates, err := service.ListExchangeRates(r.Context())
		if err != nil {
			http.Error(w, "internal error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(rates)
	}))
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Printf("listening on: %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGINT)
	fmt.Printf("got signal: %+v\n", <-quitChan)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("quit...")
}
