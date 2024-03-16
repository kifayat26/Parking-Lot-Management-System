package httpserver

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"net/http"
	"parkingManagementSystem/state"
)

func Serve(s *state.State) {
	router := chi.NewRouter()
	// middlewares
	router.Use(
		middleware.RequestID,
		httplog.RequestLogger(log.Logger),
	)

	router.Post("/pms/createUser", handleCreateUser(s))
	router.Post("/pms/createCar", handleCreateCar(s))

	router.Post("/parkCar", handleParkCar(s))
	router.Post("/unparkCar", handleUnparkCar(s))

	router.Post("/createParking", handleCreateParkingLot(s))
	router.Post("/parking-slots/maintenance", handlePutParkingSlotInMaintenance(s))
	router.Post("/parking-slots/out-of-maintenance", handlePutParkingSlotOutOfMaintenance(s))
	router.Get("/parking-lot/status", handleGetParkingLotStatus(s))
	router.Get("/history", handleGetHistoryForDay(s))

	log.Info().
		Int("port", s.Cfg.ApplicationPort).
		Msg("starting http server")

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.Cfg.ApplicationPort), router)
	if err != nil {
		log.Fatal().Err(err).Msg("http.ListenAndServe err")
	}
}
