package repository

import (
	"fmt"
	"parkingManagementSystem/models"
	"time"
)

func (repo *PgRepository) ParkCar(parkingLotID uint, carID uint) error {
	// Get the first available parking slot
	parkingSlot, err := repo.GetFirstAvailableParkingSlot(parkingLotID)
	if err != nil {
		return err
	}

	// Check if parking slot is available
	if parkingSlot == nil {
		return fmt.Errorf("no available parking slots in the specified parking lot")
	}

	// Update car data with parking slot ID
	if err := repo.DB.Model(&models.Car{}).
		Where("id = ?", carID).
		Update("parking_slot_id", parkingSlot.ID).
		Error; err != nil {
		return err
	}

	// Update parking slot data
	currentTime := time.Now()
	parkingSlot.CarID = &carID
	parkingSlot.IsBooked = true
	parkingSlot.ParkedAt = &currentTime
	parkingSlot.UnparkedAt = nil // Set unparked_at to null
	if err := repo.DB.Save(parkingSlot).Error; err != nil {
		return err
	}

	return nil
}
