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
)

type App struct {
	Db                 *db.Db
	Quit               chan bool
	WriterSerialRead   chan string
	WriterSerialWrite  chan string
	MonitorSerialRead  chan string
	MonitorSerialWrite chan string
	RestSerialRead     chan string
	RestSerialWrite    chan string

	WorkersWg *sync.WaitGroup

	Router *mux.Router
}

func (a *App) Initialize() {

	a.Quit = make(chan bool, 1)
	a.WriterSerialRead = make(chan string)
	a.WriterSerialWrite = make(chan string)
	a.MonitorSerialRead = make(chan string)
	a.MonitorSerialWrite = make(chan string)
	a.RestSerialRead = make(chan string)
	a.RestSerialWrite = make(chan string)

	a.WorkersWg = &sync.WaitGroup{}

	a.Router = mux.NewRouter()

	a.setRoutes()

	a.Db = &db.Db{}
	a.Db.Open()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(sigs, os.Interrupt)

	go func() {
		sig := <-sigs
		log.Print(sig)
		close(a.Quit)
	}()

}

func (a *App) Run() {
	done := false
	for !done {
		select {
		case <-a.Quit:
			log.Print("Main thread exit...")
			done = true
		}
	}
	a.WorkersWg.Wait()
	a.Db.Close()
	log.Print("Done.")
}

func (a *App) setRoutes() {
	a.Get("/volt", a.handleRequestLive(handler.GetVolt))
	a.Get("/hist", a.handleRequest(handler.GetHist))
	a.Get("/hist/{pg:[0-9]+}", a.handleRequest(handler.GetHist))
	a.Get("/hist/{pg:[0-9]+}/{sz:[0-9]+}", a.handleRequest(handler.GetHist))
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

type RequestHandlerFunction func(a *db.Db, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.setHeaders(w)
		handler(a.Db, w, r)
	}
}

type RequestHandlerFunctionLive func(SerialRead <-chan string, SerialWrite chan<- string, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestLive(handler RequestHandlerFunctionLive) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.setHeaders(w)
		handler(a.RestSerialRead, a.RestSerialWrite, w, r)
	}
}

func (a *App) setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type")
	//w.Header().Set("Access-Control-Allow-Credentials", "true");
}
