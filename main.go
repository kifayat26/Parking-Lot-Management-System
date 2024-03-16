package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"parkingManagementSystem/config"
	"parkingManagementSystem/httpserver"
	"parkingManagementSystem/state"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}

	appState := state.NewState(cfg)
	httpserver.Serve(appState)
}
