package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/22arw/lorafication/cmd/loraficationd/entity"
	"github.com/22arw/lorafication/internal/platform/web"
)

// CreateEntityRequest is the type that represents the request body for *Server.CreateEntity.
type CreateEntityRequest struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	SMS   *int    `json:"sms"`
}

// CreateEntityResponse is the type that represents the response body for *Server.CreateEntity.
type CreateEntityResponse struct {
	Name  string  `json:"name"`
	Email *string `json:"email"`
	SMS   *int    `json:"sms"`
}

// CreateEntity creates an entity on the lorafication server.
func (s *Server) CreateEntity(w http.ResponseWriter, r *http.Request) {
	var reqData CreateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("decode request body: %w", err))
		return
	}

	e, err := entity.CreateEntity(r.Context(), s.dbc, reqData.Name, reqData.Email, reqData.SMS)
	if err != nil {
		web.RespondError(w, r, s.logger, http.StatusInternalServerError, fmt.Errorf("create entity: %w", err))
		return
	}

	resData := CreateEntityResponse{
		Name:  e.Name,
		Email: e.Email,
		SMS:   e.SMS,
	}
	web.Respond(w, r, s.logger, http.StatusCreated, resData)
}
