package service

import (
	"schedule/pkg/model"
	"schedule/pkg/repository"
)

type SlotService struct {
	Repo *repository.SlotRepository
}

// Create a new slot
func (s *SlotService) Create(slot *model.Slot) error {
	return s.Repo.Create(slot)
}

// Get all slots for a doctor
func (s *SlotService) GetByDoctorID(doctorID uint) ([]*model.Slot, error) {
	return s.Repo.FindByDoctorID(doctorID)
}

// Get a slot by ID
func (s *SlotService) GetByID(id uint) (*model.Slot, error) {
	return s.Repo.FindByID(id)
}

// Update a slot
func (s *SlotService) Update(slot *model.Slot) error {
	return s.Repo.Update(slot)
}

// Delete a slot
func (s *SlotService) Delete(id uint) error {
	return s.Repo.Delete(id)
}
