package restless

import (
	"fmt"
	"net/http"
	"strings"
)

// Mux is a mux
type Mux struct {
	routes methodCollection
}

type routeCollection map[string]http.Handler

type methodCollection map[string]routeCollection

// NewMux creates a new Mux
func NewMux() *Mux {
	mcol := methodCollection{
		"HEAD":   routeCollection{},
		"OPTION": routeCollection{},
		"GET":    routeCollection{},
		"POST":   routeCollection{},
		"PUT":    routeCollection{},
		"DELETE": routeCollection{},
	}

	return &Mux{mcol}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fn, ok := m.routes[r.Method][r.RequestURI]; ok {
		go fn.ServeHTTP(w, r)
	} else {
		go http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// Handle adds a route handler
func (m *Mux) Handle(method, path string, fn http.Handler) error {
	return m.HandleFunc(method, path, fn.ServeHTTP)
}

// HandleFunc adds a function handler
func (m *Mux) HandleFunc(method, path string, fn http.HandlerFunc) error {
	verb := strings.ToUpper(method)

	rl, ok := m.routes[verb]
	if !ok {
		return fmt.Errorf("%s unsupported", method)
	}

	if _, ok := rl[path]; ok {
		return fmt.Errorf("%s %s is already defined", method, path)
	}

	rl[path] = fn

	return nil
}
