package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/22arw/lorafication/cmd/loraficationd/contract"
	"github.com/22arw/lorafication/internal/platform/web"
)

// CreateContractRequest is the type that represents the request body for *Server.CreateContract.
type CreateContractRequest struct {
	NodePublicKey string `json:"nodePublicKey"`
	EntityID      int    `json:"entityID"`
}

// CreateContract creates a contract between an entity and a node on the lorafication server.
func (s *Server) CreateContract(w http.ResponseWriter, r *http.Request) {
	var reqData CreateContractRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("decode request body: %w", err))
		return
	}

	if err := contract.CreateContract(r.Context(), s.dbc, reqData.NodePublicKey, reqData.EntityID); err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("create entity: %w", err))
		return
	}

	web.Respond(w, r, s.logger, http.StatusCreated, nil)
}
