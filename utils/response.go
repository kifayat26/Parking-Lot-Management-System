package utils

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
)

type CommonResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Lang    string      `json:"lang"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithError(w http.ResponseWriter, message string, statusCode int, logger zerolog.Logger) {
	//logger.Error().Str("code", "error").Str("message", message).Msg("Responding with error")
	RespondWithJSON(w, statusCode, CommonResponse{
		Code:    "error",
		Message: message,
	}, logger)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}, logger zerolog.Logger) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error().Str("code", "error").Str("message", "Failed to encode JSON response").Msg("Error encoding JSON response")
	}
}
