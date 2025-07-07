package model

import "time"

type Slot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DoctorID  uint      `json:"doctor_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsBooked  bool      `json:"is_booked"`
}
