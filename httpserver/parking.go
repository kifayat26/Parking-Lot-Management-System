package httpserver

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"net/http"
	"parkingManagementSystem/models"
	_ "parkingManagementSystem/models"
	"parkingManagementSystem/state"
	"parkingManagementSystem/utils"
	"strconv"
	"time"
)

type ReqBody struct {
	Location string `json:"location"`
	Slots    int    `json:"slots"`
}

func handleCreateParkingLot(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := s.Repository
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleCreateParkingLot").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		var reqBody ReqBody

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			logger.Error().Err(err).Msg("Failed to decode request body")
			utils.RespondWithError(w, "Failed to decode request body", http.StatusBadRequest, logger)
			return
		}

		if reqBody.Slots < 1 {
			logger.Error().Msg("Number of slots must be greater than 0")
			utils.RespondWithError(w, "Number of slots must be greater than 0", http.StatusBadRequest, logger)
			return
		}

		// Create the parking lot
		parkingLot := models.ParkingLot{
			Location: reqBody.Location,
		}
		if err := repo.Create(&parkingLot).Error; err != nil {
			logger.Error().Err(err).Msg("Failed to create parking lot")
			utils.RespondWithError(w, "Failed to create parking lot", http.StatusInternalServerError, logger)
			return
		}

		// Create parking slots with relative IDs
		for i := 1; i <= reqBody.Slots; i++ {
			parkingSlot := models.ParkingSlot{
				ParkingLotID: parkingLot.ID,
				RelativeID:   uint(i),
			}
			if err := repo.Create(&parkingSlot).Error; err != nil {
				// Rollback the created parking lot if any error occurs while creating parking slots
				repo.Delete(&parkingLot)
				logger.Error().Err(err).Msg("Failed to create parking slot")
				utils.RespondWithError(w, "Failed to create parking slot", http.StatusInternalServerError, logger)
				return
			}
		}

		// Respond with the newly created parking lot
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Parking lot created successfully",
			Data:    parkingLot,
		}, logger)
	}
}

func handleParkCar(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleParkCar").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		// Parse request parameters
		parkingLotID, err := strconv.ParseUint(r.URL.Query().Get("parking_lot_id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid parking lot ID", http.StatusBadRequest, logger)
			return
		}

		carID, err := strconv.ParseUint(r.URL.Query().Get("car_id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid car ID", http.StatusBadRequest, logger)
			return
		}

		// Check if the car is already parked
		var car models.Car
		if err := s.Repository.Find(&car, "id = ?", carID).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch car details", http.StatusInternalServerError, logger)
			return
		}
		if car.ParkingSlotID != nil {
			utils.RespondWithError(w, "Car is already parked", http.StatusBadRequest, logger)
			return
		}

		// Park the car using ParkCar function
		err = s.Repository.ParkCar(uint(parkingLotID), uint(carID))
		if err != nil {
			utils.RespondWithError(w, "Failed to park the car", http.StatusInternalServerError, logger)
			return
		}

		// Log the successful parking
		logger.Info().Str("parking_lot_id", r.URL.Query().Get("parking_lot_id")).Str("car_id", r.URL.Query().Get("car_id")).Msg("Car parked successfully")

		// Respond with success message
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Car parked successfully",
			Data:    nil,
		}, logger)
	}
}

func handleUnparkCar(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleUnparkCar").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		// Parse request parameters
		carID, err := strconv.ParseUint(r.URL.Query().Get("car_id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid car ID", http.StatusBadRequest, logger)
			return
		}

		// Find the car by ID
		var car models.Car
		if err := s.Repository.Find(&car, "id = ?", carID).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch car details", http.StatusInternalServerError, logger)
			return
		}

		// Check if the car is already unparked
		if car.ParkingSlotID == nil {
			utils.RespondWithError(w, "Car is already unparked", http.StatusBadRequest, logger)
			return
		}

		// Get the parking slot details
		var parkingSlot models.ParkingSlot
		if err := s.Repository.Find(&parkingSlot, "id = ?", *car.ParkingSlotID).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch parking slot details", http.StatusInternalServerError, logger)
			return
		}

		// Calculate parking duration and amount to be paid
		parkedAt := *parkingSlot.ParkedAt
		unparkedAt := time.Now()
		totalParkingTime := int(math.Ceil(unparkedAt.Sub(parkedAt).Hours()))
		totalAmountToBePaid := totalParkingTime * 10

		// Construct JSON data
		parkingDetails := struct {
			TotalParkingTime    int `json:"total_parking_time"`
			TotalAmountToBePaid int `json:"total_amount_to_be_paid"`
		}{
			TotalParkingTime:    totalParkingTime,
			TotalAmountToBePaid: totalAmountToBePaid,
		}

		// Update the car model
		car.ParkingSlotID = nil
		if err := s.Repository.Save(&car).Error; err != nil {
			utils.RespondWithError(w, "Failed to update car details", http.StatusInternalServerError, logger)
			return
		}

		// Update the parking slot model
		parkingSlot.IsBooked = false
		parkingSlot.CarID = nil
		parkingSlot.ParkedAt = nil
		if err := s.Repository.Save(&parkingSlot).Error; err != nil {
			utils.RespondWithError(w, "Failed to update parking slot details", http.StatusInternalServerError, logger)
			return
		}

		// Update the parking history model
		date := time.Now().Truncate(24 * time.Hour)
		var parkingHistory models.ParkingHistory
		if err := s.Repository.FirstOrCreate(&parkingHistory, "date = ?", date).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch or create parking history", http.StatusInternalServerError, logger)
			return
		}
		parkingHistory.CarsParked += 1
		parkingHistory.TotalParkingTime += totalParkingTime
		parkingHistory.TotalRevenueEarned += int64(totalAmountToBePaid)
		if err := s.Repository.Save(&parkingHistory).Error; err != nil {
			utils.RespondWithError(w, "Failed to update parking history", http.StatusInternalServerError, logger)
			return
		}

		// Log the successful unparking
		logger.Info().Str("car_id", r.URL.Query().Get("car_id")).Msg("Car unparked successfully")

		// Respond with success message and parking details
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Car unparked successfully",
			Data:    parkingDetails,
		}, logger)
	}
}

func handlePutParkingSlotInMaintenance(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handlePutParkingSlotInMaintenance").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		// Parse request parameters
		parkingSlotID, err := strconv.ParseUint(r.URL.Query().Get("parking_slot_id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid parking slot ID", http.StatusBadRequest, logger)
			return
		}

		// Fetch the parking slot details
		var parkingSlot models.ParkingSlot
		if err := s.Repository.Find(&parkingSlot, "id = ?", parkingSlotID).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch parking slot details", http.StatusInternalServerError, logger)
			return
		}

		// Check if the parking slot is already booked
		if parkingSlot.IsBooked && !parkingSlot.IsInMaintenance {
			utils.RespondWithError(w, "Parking slot is already booked", http.StatusBadRequest, logger)
			return
		}

		// Check if the parking slot is in maintenance
		if parkingSlot.IsInMaintenance {
			utils.RespondWithError(w, "Parking slot is already in maintenance", http.StatusBadRequest, logger)
			return
		}

		// Put the parking slot in maintenance
		parkingSlot.IsBooked = true
		parkingSlot.IsInMaintenance = true

		if err := s.Repository.Save(&parkingSlot).Error; err != nil {
			utils.RespondWithError(w, "Failed to put parking slot in maintenance", http.StatusInternalServerError, logger)
			return
		}

		// Log the successful putting of the parking slot in maintenance
		logger.Info().Uint("parking_slot_id", uint(parkingSlotID)).Msg("Parking slot has been put in maintenance successfully")

		// Respond with success message
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Parking slot has been put in maintenance successfully",
			Data:    nil,
		}, logger)
	}
}

func handlePutParkingSlotOutOfMaintenance(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handlePutParkingSlotOutOfMaintenance").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		// Parse request parameters
		parkingSlotID, err := strconv.ParseUint(r.URL.Query().Get("parking_slot_id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid parking slot ID", http.StatusBadRequest, logger)
			return
		}

		// Fetch the parking slot details
		var parkingSlot models.ParkingSlot
		if err := s.Repository.Find(&parkingSlot, "id = ?", parkingSlotID).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch parking slot details", http.StatusInternalServerError, logger)
			return
		}

		// Check if the parking slot is already booked
		if parkingSlot.IsBooked && !parkingSlot.IsInMaintenance {
			utils.RespondWithError(w, "Parking slot is already booked", http.StatusBadRequest, logger)
			return
		}

		// Check if the parking slot is not in maintenance
		if !parkingSlot.IsInMaintenance {
			utils.RespondWithError(w, "Parking slot is not in maintenance", http.StatusBadRequest, logger)
			return
		}

		// Put the parking slot out of maintenance
		parkingSlot.IsBooked = false
		parkingSlot.IsInMaintenance = false

		if err := s.Repository.Save(&parkingSlot).Error; err != nil {
			utils.RespondWithError(w, "Failed to put parking slot out of maintenance", http.StatusInternalServerError, logger)
			return
		}

		// Log the successful putting of the parking slot out of maintenance
		logger.Info().Uint("parking_slot_id", uint(parkingSlotID)).Msg("Parking slot has been put out of maintenance successfully")

		// Respond with success message
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Parking slot has been put out of maintenance successfully",
			Data:    nil,
		}, logger)
	}
}

func handleGetParkingLotStatus(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleGetParkingLotStatus").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		// Parse request parameters
		parkingLotID, err := strconv.ParseUint(r.URL.Query().Get("parking_lot_id"), 10, 64)
		if err != nil {
			utils.RespondWithError(w, "Invalid parking lot ID", http.StatusBadRequest, logger)
			return
		}

		// Fetch all parking slots of the parking lot
		var parkingSlots []models.ParkingSlot
		if err := s.Repository.Find(&parkingSlots, "parking_lot_id = ?", parkingLotID).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch parking slots", http.StatusInternalServerError, logger)
			return
		}

		// Create a data structure to store parking slot statuses
		var parkingLotStatus []map[string]interface{}

		// Iterate over the parking slots to collect status
		for _, slot := range parkingSlots {
			slotStatus := map[string]interface{}{
				"relative_id":       slot.RelativeID,
				"is_in_maintenance": slot.IsInMaintenance,
				"is_booked":         slot.IsBooked,
			}
			if slot.IsBooked {
				slotStatus["carID"] = *slot.CarID
			}
			parkingLotStatus = append(parkingLotStatus, slotStatus)
		}

		// Log the successful fetching of parking slot statuses
		logger.Info().Uint64("parking_lot_id", parkingLotID).Msg("Parking slot statuses fetched successfully")

		// Respond with parking slot statuses
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "Parking slot statuses fetched successfully",
			Data:    parkingLotStatus,
		}, logger)
	}
}

func handleGetHistoryForDay(s *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := log.With().
			Str("handler", "handleGetHistoryForDay").
			Str("request_id", middleware.GetReqID(ctx)).
			Logger()

		// Parse the request parameters
		day := r.URL.Query().Get("day")
		if day == "" {
			utils.RespondWithError(w, "Day parameter is required", http.StatusBadRequest, logger)
			return
		}

		// Parse the day parameter to time.Time
		date, err := time.Parse("2006-01-02", day)
		if err != nil {
			utils.RespondWithError(w, "Invalid date format. Please provide the date in YYYY-MM-DD format", http.StatusBadRequest, logger)
			return
		}

		// Fetch history data for the specified day
		var history models.ParkingHistory
		if err := s.Repository.Find(&history, "date = ?", date).Error; err != nil {
			utils.RespondWithError(w, "Failed to fetch history data", http.StatusInternalServerError, logger)
			return
		}

		// Log the successful fetching of history data
		logger.Info().Time("date", date).Msg("History data fetched successfully")

		// Respond with history data
		utils.RespondWithJSON(w, http.StatusOK, utils.CommonResponse{
			Code:    "success",
			Message: "History data fetched successfully",
			Data:    history,
		}, logger)
	}
}
