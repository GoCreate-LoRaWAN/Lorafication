// Package server contains all of the handlers for the lorafication API.
package server

import (
	"errors"
	"net/http"
	"runtime"

	"github.com/george-e-shaw-iv/lorafication/cmd/loraficationd/config"
	"github.com/george-e-shaw-iv/lorafication/internal/platform/web"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Server implements http.Handler and handles incoming HTTP connections.
type Server struct {
	Config *config.Config
	Logger *zap.Logger

	http.Handler
}

// NewServer returns a reference to a Server type with the fields and handlers properly
// set.
func NewServer(cfg *config.Config, logger *zap.Logger) *Server {
	s := Server{
		Config: cfg,
		Logger: logger,
	}

	r := httprouter.New()

	// Boilerplate Routes
	s.boilerplate(r)

	// Notification Routes
	r.HandlerFunc(http.MethodPost, "/notify", s.Notify)

	// Wrap handler in middleware that handles logging and verification of the
	// RequestID.
	s.Handler = web.RequestMW(logger, r)

	return &s
}

// boilerplate initializes boilerplate routes on the given router for things like panic handling,
// kubernetes probes, etc.
func (s *Server) boilerplate(r *httprouter.Router) {
	// Gracefully handle and recover from web panics.
	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]

		s.Logger.Error("captured http panic",
			zap.Reflect("panic", i), zap.String("stack", string(stack)))

		web.RespondError(w, r, s.Logger, http.StatusInternalServerError,
			errors.New(http.StatusText(http.StatusInternalServerError)))
	}

	// Not found (404) handler being overridden for logging purposes. This helps distinguish 404s
	// due to resources not being found from unregistered routes being reached in terms of logging.
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Error("unregistered route attempting to be reached",
			zap.String("requestURI", r.RequestURI))

		web.RespondError(w, r, s.Logger, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
	})

	// probeHandler for kubernetes probes. Right now the logic is non-existent and indistinguishable
	// for the ready and healthy handlers, but eventually there is a chance it will not be.
	probeHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Routes for Kubernetes probes.
	r.HandlerFunc(http.MethodGet, "/ready", probeHandler)
	r.HandlerFunc(http.MethodGet, "/healthy", probeHandler)
}
