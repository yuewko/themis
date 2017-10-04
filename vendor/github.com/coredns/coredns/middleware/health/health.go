// Package health implements an HTTP handler that responds to health checks.
package health

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

var once sync.Once

type health struct {
	Addr string

	ln  net.Listener
	mux *http.ServeMux

	// A slice of Healthers that the health middleware will poll every second for their health status.
	h []Healther
	sync.RWMutex
	ok bool // ok is the global boolean indicating an all healthy middleware stack
}

func (h *health) Startup() error {
	if h.Addr == "" {
		h.Addr = defAddr
	}

	once.Do(func() {
		ln, err := net.Listen("tcp", h.Addr)
		if err != nil {
			log.Printf("[ERROR] Failed to start health handler: %s", err)
			return
		}

		h.ln = ln

		h.mux = http.NewServeMux()

		h.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if h.Ok() {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, ok)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
		})

		go func() {
			http.Serve(h.ln, h.mux)
		}()
	})
	return nil
}

func (h *health) Shutdown() error {
	if h.ln != nil {
		return h.ln.Close()
	}
	return nil
}

const (
	ok      = "OK"
	defAddr = ":8080"
	path    = "/health"
)
