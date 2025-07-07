package service

import (
	"appointment/pkg/model"
	"appointment/pkg/repository"
)

type AppointmentService struct {
	Repo *repository.AppointmentRepository
}

// Create a new appointment
func (s *AppointmentService) Create(app *model.Appointment) error {
	return s.Repo.Create(app)
}

// Get appointment by ID
func (s *AppointmentService) GetByID(id uint) (*model.Appointment, error) {
	return s.Repo.FindByID(id)
}

// Get all appointments for a patient
func (s *AppointmentService) GetByPatientID(patientID uint) ([]*model.Appointment, error) {
	return s.Repo.FindByPatientID(patientID)
}

// Get all appointments for a doctor
func (s *AppointmentService) GetByDoctorID(doctorID uint) ([]*model.Appointment, error) {
	return s.Repo.FindByDoctorID(doctorID)
}

// Get appointment by SlotID
func (s *AppointmentService) GetBySlotID(slotID uint) (*model.Appointment, error) {
	return s.Repo.FindBySlotID(slotID)
}

// Update an appointment
func (s *AppointmentService) Update(app *model.Appointment) error {
	return s.Repo.Update(app)
}

// Delete an appointment
func (s *AppointmentService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
