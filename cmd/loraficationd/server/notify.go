package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/22arw/lorafication/cmd/loraficationd/nodes"
	"github.com/22arw/lorafication/cmd/loraficationd/notifications"
	"github.com/22arw/lorafication/internal/platform/web"
)

// NotifyRequest contains the structure of the request data for the notify route.
type NotifyRequest struct {
	NodeID  int    `json:"nodeID"`
	Message string `json:"message"`
}

// Notify handles the /notify endpoint to notify all entities subscribed to a node.
func (s *Server) Notify(w http.ResponseWriter, r *http.Request) {
	var reqData NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("decode request body: %w", err))
		return
	}

	node, err := nodes.ResolveNode(s.dbc, reqData.NodeID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
		}

		web.RespondError(w, r, s.logger, statusCode, fmt.Errorf("resolve node from id: %w", err))
		return
	}

	contracts, err := notifications.ResolveContracts(s.dbc, reqData.NodeID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
		}

		web.RespondError(w, r, s.logger, statusCode, fmt.Errorf("resolve contracts from node id: %w", err))
		return
	}

	subject := fmt.Sprintf("LoRafication: Notification from %s Node", node.Name)
	for i := range contracts {
		if contracts[i].Email != nil {
			if err := s.mailer.Send(*contracts[i].Email, subject, reqData.Message); err != nil {
				// TODO: This should have a failsafe mechanism
				web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("send notification: %w", err))
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
