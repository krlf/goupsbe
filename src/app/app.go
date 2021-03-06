package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"upsbe/config"
	"upsbe/db"
	"upsbe/handler"
	"upsbe/types"

	"github.com/gorilla/mux"
)

type App struct {
	db           *db.Db
	quit         chan bool
	serialStream chan types.StringStream
	workersWg    *sync.WaitGroup
	router       *mux.Router
	upsConfig    *config.Config
}

func (a *App) Initialize(upsConfig *config.Config, db *db.Db) {

	a.quit = make(chan bool, 1)

	a.serialStream = make(chan types.StringStream)

	a.workersWg = &sync.WaitGroup{}

	a.router = mux.NewRouter()

	a.setRoutes()

	a.router.Use(corsMethodMiddleware(a.router))

	a.upsConfig = upsConfig
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

func (a *App) WorkersWgGet() *sync.WaitGroup {
	return a.workersWg
}

func (a *App) SerialStreamGet() chan types.StringStream {
	return a.serialStream
}

func (a *App) QuitFlagGet() chan bool {
	return a.quit
}

func (a *App) RouterGet() *mux.Router {
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

	a.Get("/config", a.handleRequestConfig(handler.GetConfig))
	a.Put("/config", a.handleRequestConfig(handler.SetConfig))
}

func (a *App) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, handler).Methods(http.MethodGet, http.MethodOptions)
}

func (a *App) Put(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	a.router.HandleFunc(path, handler).Methods(http.MethodPut, http.MethodOptions)
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func corsMethodMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var allMethods []string

			err := r.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
				if route.Match(req, &mux.RouteMatch{}) {
					methods, err := route.GetMethods()
					if err != nil {
						return err
					}
					allMethods = append(allMethods, methods...)
				}
				return nil
			})

			if err == nil {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(unique(allMethods), ","))
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,content-type")
				// w.Header().Set("Access-Control-Max-Age", "86400")
				// w.Header().Set("Access-Control-Allow-Credentials", "true");

				if req.Method == "OPTIONS" {
					return
				}
			}

			next.ServeHTTP(w, req)
		})
	}
}

type RequestHandlerFunctionDb func(db *db.Db, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestDb(handler RequestHandlerFunctionDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.db, w, r)
	}
}

type RequestHandlerFunctionStream func(stream types.StringStream, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestStream(handler RequestHandlerFunctionStream) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stream := types.StringStreamCreate()
		a.serialStream <- stream
		handler(stream, w, r)
		close(stream.Write)
	}
}

type RequestHandlerFunctionConfig func(upsConfig *config.Config, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequestConfig(handler RequestHandlerFunctionConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.upsConfig, w, r)
	}
}
