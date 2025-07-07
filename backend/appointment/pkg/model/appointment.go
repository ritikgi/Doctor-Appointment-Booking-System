package model

import "time"

type Appointment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PatientID uint      `json:"patient_id"`
	DoctorID  uint      `json:"doctor_id"`
	SlotID    uint      `json:"slot_id"`
	Status    string    `json:"status"` // Upcoming, Completed, Cancelled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
