package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"

	"github.com/3n0ugh/allotropes/framework/swagger"
	"github.com/3n0ugh/allotropes/internal/errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	for _, c := range a.Controllers {
		for _, r := range c.Routes {
			router.Method(r.Method, r.Path, r.Handler)
		}
	}

	if a.SwaggerEnabled {
		for _, c := range a.Controllers {
			swagger.S.AddTag(c.Name, c.Description)
			for _, r := range c.Routes {
				swagger.S.SetPaths(r, c.Name)
				if _, ok := swagger.S.Components.Schemas[reflect.TypeOf(r.Request).Name()]; !ok {
					swagger.S.Components.Schemas[reflect.TypeOf(r.Request).Name()] = swagger.ComponentSchema{
						XSwaggerRouterModel: swagger.RouterModelPrefix + reflect.TypeOf(r.Request).Name(),
					}
				}

				if _, ok := swagger.S.Components.Schemas[reflect.TypeOf(r.Response).Name()]; !ok {
					swagger.S.Components.Schemas[reflect.TypeOf(r.Response).Name()] = swagger.ComponentSchema{
						XSwaggerRouterModel: swagger.RouterModelPrefix + reflect.TypeOf(r.Response).Name(),
					}
				}
			}
		}

		swagger.S.SetSchema(reflect.TypeOf(errors.Error{}))

		if err := swagger.Init(a.Name); err != nil {
			log.Println("swagger init: ", err)
		}

		sh := http.StripPrefix("/swagger/", http.FileServer(http.Dir("./framework/swagger/dist/")))
		router.Method(http.MethodGet, "/swagger/*", sh)
	}

	a.server = http.Server{
		Addr:    ":" + strconv.Itoa(a.Port),
		Handler: router,
	}
}
