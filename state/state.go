package state

import (
	"github.com/rs/zerolog/log"
	"parkingManagementSystem/config"
	"parkingManagementSystem/repository"
)

type State struct {
	Cfg        *config.Config
	Repository *repository.PgRepository
}

func NewState(cfg *config.Config) *State {
	db, err := repository.NewPgRepository(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("pg repository error")
	}
	err = db.Migrate()
	if err != nil {
		log.Fatal().Err(err).Msg("pg repository error")
	}

	return &State{
		Cfg:        cfg,
		Repository: db,
	}
}
