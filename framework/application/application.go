package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/3n0ugh/allotropes/framework/swagger"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

type App struct {
	Name           string
	Port           int
	Controllers    []Controller
	Config         any
	Logger         log.Logger
	Cache          any
	server         http.Server
	SwaggerEnabled bool
}

func (a *App) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Printf("Error running server: %s", err)
			close(stop)
		}
	}()

	fmt.Println("Server running on", a.server.Addr)
	<-stop
	if err := a.server.Shutdown(context.TODO()); err != nil {
		log.Println(err)
	}
}

func (a *App) Setup() {
	router := mux.NewRouter().StrictSlash(true)
	for _, c := range a.Controllers {
		for _, r := range c.Routes {
			router.Handle(r.Path, r.Handler).Methods(r.Method)
		}
	}

	if a.SwaggerEnabled {
		for _, c := range a.Controllers {
			swagger.S.AddTag(c.Name, c.Description)
			for _, r := range c.Routes {
				swagger.S.SetPaths(r, c.Name)
			}
		}

		if err := swagger.Init(a.Name); err != nil {
			log.Println("swagger init: ", err)
		}

		sh := http.StripPrefix("/swagger/", http.FileServer(http.Dir("./dist/")))
		router.PathPrefix("/swagger/").Handler(middleware.Logger(sh))
	}

	a.server = http.Server{
		Addr:    ":" + strconv.Itoa(a.Port),
		Handler: router,
	}
}
