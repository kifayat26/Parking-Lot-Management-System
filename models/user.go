package models

type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Cars []Car // One-to-Many relationship: One user can have multiple cars
}

type Car struct {
	ID            uint  `gorm:"primaryKey"`
	UserID        uint  // Foreign key to User.ID
	ParkingSlotID *uint `json:"parking_slot_id,omitempty"` // Nullable reference to ParkingSlot
}
