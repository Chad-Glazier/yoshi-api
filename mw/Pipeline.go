package mw

import (
	"database/sql"
	"net/http"
	"yoshi/db/user"
)

type CleanupFunc = func(*Resources)
type MiddlewareFunc = func(*Resources, http.ResponseWriter, *http.Request) (bool, CleanupFunc)
type EndHandler = func(*Resources, http.ResponseWriter, *http.Request)
type Resources = struct {
	DB      *sql.DB
	Session *user.Session
	Body	any
}

type Pipeline struct {
	middleware []MiddlewareFunc
	endHandler EndHandler

	// common data that may be created/consumed by various
	// middleware or the end handler.
	resources Resources

	// used to cleanup the request (close memory leaks and stuff)
	cleanupFuncs []CleanupFunc
}

func (p *Pipeline) Run(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range p.middleware {
		proceed, cleanupFunc := middleware(&p.resources, w, r)
		if cleanupFunc != nil {
			p.cleanupFuncs = append(p.cleanupFuncs, cleanupFunc)
		}
		if !proceed {
			p.cleanup()
			return
		}
	}
	p.endHandler(&p.resources, w, r)
	p.cleanup()
}

func (p *Pipeline) cleanup() {
	for _, cleanupFunc := range p.cleanupFuncs {
		cleanupFunc(&p.resources)
	}
}

func (p *Pipeline) Use(middleware ...MiddlewareFunc) {
	p.middleware = append(p.middleware, middleware...)
}

func (p *Pipeline) SetEndHandler(h EndHandler) {
	p.endHandler = h
}

func NewPipeline(endHandler EndHandler, middleware ...MiddlewareFunc) *Pipeline {
	var pipeline Pipeline
	pipeline.Use(middleware...)
	pipeline.SetEndHandler(endHandler)
	return &pipeline
}
