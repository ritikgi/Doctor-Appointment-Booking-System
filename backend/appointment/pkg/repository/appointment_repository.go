package repository

import (
	"appointment/pkg/model"

	"gorm.io/gorm"
)

type AppointmentRepository struct {
	DB *gorm.DB
}

// Create a new appointment
func (r *AppointmentRepository) Create(app *model.Appointment) error {
	return r.DB.Create(app).Error
}

// Find appointment by ID
func (r *AppointmentRepository) FindByID(id uint) (*model.Appointment, error) {
	var app model.Appointment
	err := r.DB.First(&app, id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Find all appointments for a patient
func (r *AppointmentRepository) FindByPatientID(patientID uint) ([]*model.Appointment, error) {
	var apps []*model.Appointment
	err := r.DB.Where("patient_id = ?", patientID).Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// Find all appointments for a doctor
func (r *AppointmentRepository) FindByDoctorID(doctorID uint) ([]*model.Appointment, error) {
	var apps []*model.Appointment
	err := r.DB.Where("doctor_id = ?", doctorID).Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// Find appointment by SlotID
func (r *AppointmentRepository) FindBySlotID(slotID uint) (*model.Appointment, error) {
	var app model.Appointment
	// Only consider appointments that are not cancelled
	err := r.DB.Where("slot_id = ? AND status != ?", slotID, "Cancelled").First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Update an appointment
func (r *AppointmentRepository) Update(app *model.Appointment) error {
	return r.DB.Save(app).Error
}

// Delete an appointment
func (r *AppointmentRepository) Delete(id uint) error {
	return r.DB.Delete(&model.Appointment{}, id).Error
}
