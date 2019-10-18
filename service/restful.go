package service

import (
	"fmt"
	"net/http"
)
type RestfulRouter map[string]http.HandlerFunc
//basic RestfulService entry point
type RestfulService struct {
	ctx *Context
	router RestfulRouter
}

func NewRestfulService(ctx *Context) RestfulService {
	svc := RestfulService{
		ctx: ctx,
		router: make(RestfulRouter),
	}

	svc.router["/"] = svc.restfulRoot()
	svc.router["/health"] = svc.health()

	return svc
}

func (rs RestfulService) Router() RestfulRouter {
	return rs.router
}

func (rs RestfulService) restfulRoot() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, err := fmt.Fprintln(w, "Available endpoints: ")
		if err != nil {
			rs.ctx.Logger.Errorf("failed to display available endpoints info")
		}
		for path := range rs.router {
			_, err = fmt.Fprintln(w, r.Host+path)
			if err != nil {
				rs.ctx.Logger.Errorf("failed to display available endpoints info")
			}
		}
	}
}

func (rs RestfulService) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "health check for SDK port %v : OK\n", rs.ctx.Cfg.Network.SDKAddress)
		if err != nil {
			rs.ctx.Logger.Errorf("failed to display SDK port health check info")
		}
	}
}

