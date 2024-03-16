package models

import "time"

type ParkingLot struct {
	ID       uint          `gorm:"primaryKey" json:"id"`
	Location string        `json:"location"`
	Slots    []ParkingSlot `json:"slots"`
}

type ParkingSlot struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ParkingLotID    uint       `json:"-"`
	RelativeID      uint       `json:"relative_id"`
	IsBooked        bool       `gorm:"default:false" json:"is_booked"`
	IsInMaintenance bool       `gorm:"default:false" json:"is_in_maintenance"`
	CarID           *uint      `json:"car_id,omitempty"` // Nullable reference to Car
	ParkedAt        *time.Time `json:"parked_at,omitempty"`
	UnparkedAt      *time.Time `json:"unparked_at,omitempty"`
}

type ParkingHistory struct {
	Date               time.Time `gorm:"primaryKey" json:"date"`
	CarsParked         int       `json:"cars_parked"`
	TotalParkingTime   int       `json:"total_parking_time" gorm:"default:0"`   // Total parking time in minutes
	TotalRevenueEarned int64     `json:"total_revenue_earned" gorm:"default:0"` // Total revenue earned
}
