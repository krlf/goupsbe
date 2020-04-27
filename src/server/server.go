package server

import (
	"../app"
	"../config"
	"context"
	"log"
	"net/http"
	"time"
)

func Worker(a *app.App) {

	a.WorkersWg.Add(1)
	defer a.WorkersWg.Done()

	log.Print("Server thread starting...")

	srv := &http.Server{
		Addr:    ":" + config.GetListenPort(),
		Handler: a.Router}

	go func() {
		log.Print("HTTP Server Listener thread starting...")
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal("[Server] Error ListenAndServe(): %v", err)
		}
	}()

	<-a.Quit
	log.Print("HTTP Server Listener shutdown...")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	srv.Shutdown(ctxShutDown)
	log.Print("Server thread exit...")
}
