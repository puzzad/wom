package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/csmith/envflag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/puzzad/wom"
)

var port = flag.Int("port", 3000, "Port to listen for HTTP requests")

func main() {
	envflag.Parse()
	log.Println("Wise old Man is starting...")

	if err := wom.ConnectToDatabase(); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/mail/subscribe", wom.SubscribeToMailingList)
	r.Post("/mail/confirm", wom.ConfirmMailingListSubscription)
	r.Post("/mail/unsubscribe", wom.UnsubscribeFromMailingList)
	r.Post("/mail/contact", wom.SendContactForm)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: r,
	}

	go func() {
		log.Printf("Listening on port %d...\n", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen failed: %s\n", err)
		}
	}()

	<-done
	log.Println("Signal received, shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
