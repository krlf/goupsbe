package app

import (
	"../db"
	"../handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"../types"
	"../config"
)

type App struct {
	db                 *db.Db
	quit               chan bool
	serialStream chan types.StringStream
	workersWg *sync.WaitGroup
	router *mux.Router
}

func (a *App) Initialize(c *config.Config, db *db.Db) {

	a.quit = make(chan bool, 1)

	a.serialStream = make(chan types.StringStream)

	a.workersWg = &sync.WaitGroup{}

	a.router = mux.NewRouter()

	a.setRoutes()

	a.db = db

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(sigs, os.Interrupt)

	go func() {
		sig := <-sigs
		log.Print(sig)
		close(a.quit)
	}()

}

func (a *App)WorkersWgGet() *sync.WaitGroup {
	return a.workersWg
}

func (a *App)SerialStreamGet() chan types.StringStream {
	return a.serialStream
}

func (a *App)QuitFlagGet() chan bool {
	return a.quit
}

func (a *App)RouterGet() *mux.Router {
	return a.router
}



func (a *App) Run() {
	done := false
	for !done {
		select {
		case <-a.quit:
			log.Print("Main thread exit...")
			done = true
		}
	}
	a.workersWg.Wait()
	log.Print("Done.")
}

func (a *App) setRoutes() {
	a.Get("/volt", a.handleRequestStream(handler.GetVolt))
	a.Get("/hist", a.handleRequestDb(handler.GetHist))
	a.Get("/hist/{pg:[0-9]+}", a.handleRequestDb(handler.GetHist))
	a.Get("/hist/{pg:[0-9]+}/{sz:[0-9]+}", a.handleRequestDb(handler.GetHist))
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, f).Methods("GET")
}

type RequestHandlerFunctionDb func(a *db.Db, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestDb(handler RequestHandlerFunctionDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.setHeaders(w)
		handler(a.db, w, r)
	}
}

type RequestHandlerFunctionStream func(stream types.StringStream, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestStream(handler RequestHandlerFunctionStream) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.setHeaders(w)
		stream := types.StringStreamCreate()
		a.serialStream <- stream;
		handler(stream, w, r)
		close(stream.Write)
	}
}

func (a *App) setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type")
	//w.Header().Set("Access-Control-Allow-Credentials", "true");
}
