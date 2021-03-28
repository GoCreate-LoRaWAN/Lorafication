package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/22arw/lorafication/cmd/loraficationd/node"
	"github.com/22arw/lorafication/internal/platform/web"
)

// CreateNodeRequest is the type that represents the request body for *Server.CreateNode.
type CreateNodeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateNodeResponse is the type that represents the response body for *Server.CreateNode.
type CreateNodeResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PublicKey   string `json:"publicKey"`
	Secret      string `json:"secret"`
}

// CreateNode creates a node on the lorafication server.
func (s *Server) CreateNode(w http.ResponseWriter, r *http.Request) {
	var reqData CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("decode request body: %w", err))
		return
	}

	n, err := node.CreateNode(r.Context(), s.dbc, reqData.Name, reqData.Description)
	if err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("create node: %w", err))
		return
	}

	resData := CreateNodeResponse{
		Name:        n.Name,
		Description: n.Description,
		PublicKey:   n.PublicKey,
		Secret:      n.Secret,
	}
	web.Respond(w, r, s.logger, http.StatusCreated, resData)
}
