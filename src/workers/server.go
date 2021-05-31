package workers

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"upsbe/config"

	"github.com/gorilla/mux"
)

func Server(quitFlag chan bool, workersWg *sync.WaitGroup, upsConfig *config.Config, router *mux.Router) {

	workersWg.Add(1)
	defer workersWg.Done()

	log.Print("Server thread starting...")

	srv := &http.Server{
		Addr:    ":" + upsConfig.ListenPortGet(),
		Handler: router}

	go func() {
		log.Print("HTTP Server Listener thread starting...")
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal("[Server] Error ListenAndServe(): %v", err)
		}
	}()

	<-quitFlag
	log.Print("HTTP Server Listener shutdown...")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	srv.Shutdown(ctxShutDown)
	log.Print("Server thread exit...")
}
