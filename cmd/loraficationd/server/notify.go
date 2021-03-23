package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/22arw/lorafication/cmd/loraficationd/contract"
	"github.com/22arw/lorafication/cmd/loraficationd/node"
	"github.com/22arw/lorafication/internal/platform/web"
)

// NotifyRequest is a representation of the request body for the *Server.Notify handler.
type NotifyRequest struct {
	PublicKey string `json:"publicKey"` // PublicKey corresponds to a node public key (primary key of a node).
	Secret    string `json:"secret"`    // Secret corresponds to the secret stored in the same row^.
	Message   string `json:"message"`
}

// Notify notifies all entities subscribed to a node using the provided message.
func (s *Server) Notify(w http.ResponseWriter, r *http.Request) {
	var reqData NotifyRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("decode request body: %w", err))
		return
	}

	n, err := node.AuthenticateNode(r.Context(), s.dbc, reqData.PublicKey, reqData.Secret)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusUnauthorized
		}

		web.RespondError(w, r, s.logger, statusCode, fmt.Errorf("resolve node from public key and secret: %w", err))
		return
	}

	contracts, err := contract.ResolveContracts(r.Context(), s.dbc, reqData.PublicKey)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
		}

		web.RespondError(w, r, s.logger, statusCode, fmt.Errorf("resolve contracts from node id: %w", err))
		return
	}

	subject := fmt.Sprintf("LoRafication: Notification from %s Node", n.Name)
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
