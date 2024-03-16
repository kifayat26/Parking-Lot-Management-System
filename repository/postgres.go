package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"parkingManagementSystem/models"
)

type PgRepository struct {
	*gorm.DB
}

func NewPgRepository(databaseUrl string) (*PgRepository, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PgRepository{DB: db}, nil
}

func (repo *PgRepository) Migrate() error {
	// Migrate User and Car models
	if err := repo.DB.Migrator().AutoMigrate(&models.User{}, &models.Car{}); err != nil {
		return err
	}

	// Migrate ParkingSlot model with index
	if err := repo.DB.Migrator().AutoMigrate(&models.ParkingSlot{}); err != nil {
		return err
	}

	// Define index for ParkingSlot model
	if err := repo.DB.Exec("CREATE INDEX idx_parking_slot_composite ON parking_slots (parking_lot_id, is_booked, relative_id)").Error; err != nil {
		return err
	}

	return nil
}

func (repo *PgRepository) GetFirstAvailableParkingSlot(parkingLotID uint) (*models.ParkingSlot, error) {
	var parkingSlot models.ParkingSlot
	if err := repo.DB.Where("parking_lot_id = ? AND is_booked = ?", parkingLotID, false).
		Order("relative_id").
		First(&parkingSlot).
		Error; err != nil {
		return nil, err
	}
	return &parkingSlot, nil
}
