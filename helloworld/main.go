package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gnanasuriyan/go-micro-services-http/helloworld/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api ", log.LstdFlags)

	sm := http.NewServeMux()
	sm.Handle("/", handlers.NewHello(logger))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
		logger.Println("server is running and listening port 9090")
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	logger.Println("Received terminate, graceful shutdown", sig)
	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(timeoutContext)
}
