package application

import (
	"encoding/json"
	"log"
	"net/http"
)

type Controller struct {
	Name        string
	Description string
	Routes      []Route
}

type Route struct {
	Name        string
	Description string
	Method      string
	Path        string
	Headers     map[string]string
	Middlewares []func(http.Handler) http.Handler
	Handler     HandlerFunc
	Request     any
	Response    any
}

func (r Route) GetRequestModel() any                              { return r.Request }
func (r Route) GetResponseModel() any                             { return r.Response }
func (r Route) GetMethod() string                                 { return r.Method }
func (r Route) GetDescription() string                            { return r.Description }
func (r Route) GetMiddlewares() []func(http.Handler) http.Handler { return r.Middlewares }
func (r Route) GetHeaders() map[string]string                     { return r.Headers }
func (r Route) GetPath() string                                   { return r.Path }

type HandlerFunc func(w http.ResponseWriter, r *http.Request) (response any, err error)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	res, err := h(w, r)
	if err != nil {
		log.Printf("err: %s", err.Error())
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("err: %s", err.Error())
	}
}
