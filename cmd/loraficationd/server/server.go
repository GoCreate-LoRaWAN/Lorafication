// Package server contains all of the handlers for the lorafication API.
package server

import (
	"errors"
	"net/http"
	"runtime"

	"github.com/22arw/lorafication/cmd/loraficationd/config"
	"github.com/22arw/lorafication/internal/mail"
	"github.com/22arw/lorafication/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// Server implements http.Handler and handles incoming HTTP connections.
type Server struct {
	config *config.Config
	logger *zap.Logger
	dbc    *sqlx.DB
	mailer *mail.Mailer

	http.Handler
}

// NewServer returns a reference to a Server type with the fields and handlers properly
// set.
func NewServer(cfg *config.Config, logger *zap.Logger, dbc *sqlx.DB, mailer *mail.Mailer) *Server {
	s := Server{
		config: cfg,
		logger: logger,
		dbc:    dbc,
		mailer: mailer,
	}

	r := httprouter.New()

	// Boilerplate Routes
	s.boilerplate(r)

	// Entity Routes
	r.HandlerFunc(http.MethodPost, "/entity", s.CreateEntity)

	// Node Routes
	r.HandlerFunc(http.MethodPost, "/node", s.CreateNode)

	// Node/Entity Contract Routes
	r.HandlerFunc(http.MethodPost, "/contract", s.CreateContract)

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

		s.logger.Error("captured http panic",
			zap.Reflect("panic", i), zap.String("stack", string(stack)))

		web.RespondError(w, r, s.logger, http.StatusInternalServerError,
			errors.New(http.StatusText(http.StatusInternalServerError)))
	}

	// Not found (404) handler being overridden for logging purposes. This helps distinguish 404s
	// due to resources not being found from unregistered routes being reached in terms of logging.
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Error("unregistered route attempting to be reached",
			zap.String("requestURI", r.RequestURI))

		web.RespondError(w, r, s.logger, http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)))
	})

	// probeHandler is for kubernetes probes.
	probeHandler := func(w http.ResponseWriter, r *http.Request) {
		if err := s.dbc.Ping(); err == nil {

			// Ping by itself is un-reliable, the connections are cached. This
			// ensures that the database is still running by executing a harmless
			// dummy query against it.
			if _, err = s.dbc.Exec("SELECT true"); err == nil {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
	}

	// Routes for Kubernetes probes.
	r.HandlerFunc(http.MethodGet, "/ready", probeHandler)
	r.HandlerFunc(http.MethodGet, "/healthy", probeHandler)
}
