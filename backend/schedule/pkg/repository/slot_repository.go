package repository

import (
	"schedule/pkg/model"

	"gorm.io/gorm"
)

type SlotRepository struct {
	DB *gorm.DB
}

// Create a new slot
func (r *SlotRepository) Create(slot *model.Slot) error {
	return r.DB.Create(slot).Error
}

// Find all slots for a doctor
func (r *SlotRepository) FindByDoctorID(doctorID uint) ([]*model.Slot, error) {
	var slots []*model.Slot
	err := r.DB.Where("doctor_id = ?", doctorID).Find(&slots).Error
	if err != nil {
		return nil, err
	}
	return slots, nil
}

// Find a slot by ID
func (r *SlotRepository) FindByID(id uint) (*model.Slot, error) {
	var slot model.Slot
	err := r.DB.First(&slot, id).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

// Update a slot
func (r *SlotRepository) Update(slot *model.Slot) error {
	return r.DB.Save(slot).Error
}

// Delete a slot
func (r *SlotRepository) Delete(id uint) error {
	return r.DB.Delete(&model.Slot{}, id).Error
}
