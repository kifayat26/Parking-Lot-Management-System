package httpserver

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"parkingManagementSystem/models"
	"parkingManagementSystem/state"
	"parkingManagementSystem/utils"
	"strconv"
)

func handleCreateUser(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleCreateUser").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Error().Err(err).Msg("Failed to decode request body")
			utils.RespondWithError(w, "Failed to decode request body", http.StatusBadRequest, logger)
			return
		}

		// Create user using PgRepository
		if err := s.Repository.Create(&user).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to create user")
			utils.RespondWithError(w, "Failed to create user", http.StatusInternalServerError, logger)
			return
		}

		// Encode response with created user
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "User created successfully",
			Data:    user,
		}, logger)
	}
}

func handleCreateCar(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleCreateCar").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()
		// Get userID from URL path
		userID := chi.URLParam(r, "userID")
		uid, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid user ID", http.StatusBadRequest, logger)
			return
		}

		// Decode request body to Car struct
		var car models.Car

		if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
			logger.Error().Err(err).Msg("Failed to decode request body")
			utils.RespondWithError(w, "Failed to decode request body", http.StatusBadRequest, logger)
			return
		}

		// Set userID to the car
		car.UserID = uint(uid)

		// Create car using PgRepository
		if err := s.Repository.Create(&car).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to create car")
			utils.RespondWithError(w, "Failed to create car", http.StatusInternalServerError, logger)
			return
		}

		// Encode response with created car
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Car created successfully",
			Data:    car,
		}, logger)
	}
}
