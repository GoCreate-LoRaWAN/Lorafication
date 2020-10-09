// Package web contains various server handler helpers, such as global
// request middleware.
package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// Response is the format used for all the responses.
type Response struct {
	Results interface{}     `json:"results"`
	Errors  []ResponseError `json:"errors,omitempty"`
}

// ResponseError is the format used for response errors.
type ResponseError struct {
	Message string `json:"message"`
}

// Error implements the error interface.
func (a ResponseError) Error() string {
	return a.Message
}

// Respond sends a response with a status code.
func Respond(w http.ResponseWriter, r *http.Request, logger *zap.Logger, code int, data interface{}, errs ...error) {
	var respErrs []ResponseError

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error("error in successful request", zap.Error(err))
			respErrs = append(respErrs, ResponseError{Message: err.Error()})
		}
	}

	resp := Response{
		Results: data,
		Errors:  respErrs,
	}

	writeResponse(w, r, logger, code, &resp)
}

// RespondError sends an error response with a status code. The error is automatically logged for you.
// If the error implements StatusCoder, the provided status code will be used.
func RespondError(w http.ResponseWriter, r *http.Request, logger *zap.Logger, code int, err error) {
	logger.Error("error in unsuccessful request", zap.Error(err))

	if code >= http.StatusInternalServerError {

		// Respond with generic error. Error messages and and codes may potentially contain
		// sensitive information or help an attacker.
		code = http.StatusInternalServerError
		err = errors.New(http.StatusText(http.StatusInternalServerError))
	}

	resp := Response{
		Errors: []ResponseError{
			{
				Message: err.Error(),
			},
		},
	}

	writeResponse(w, r, logger, code, &resp)
}

// writeResponse marshals the response to json and writes it to the response writer.
func writeResponse(w http.ResponseWriter, r *http.Request, logger *zap.Logger, code int, resp *Response) {
	if code == http.StatusNoContent || resp == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		RespondError(w, r, logger, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(b); err != nil {
		logger.Error("write response body", zap.Error(err))
	}
}
