package handler

import (
	"appointment/pkg/model"
	"appointment/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	Service *service.AppointmentService
}

// AppointmentRequest represents the request body for creating an appointment
// swagger:model
type AppointmentRequest struct {
	DoctorID uint `json:"doctor_id"`
	SlotID   uint `json:"slot_id"`
}

// CreateAppointment godoc
// @Summary Book a new appointment
// @Description Patient books an appointment for a slot with a doctor
// @Tags appointments
// @Accept json
// @Produce json
// @Param appointment body AppointmentRequest true "Appointment info"
// @Success 201 {object} model.Appointment
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /appointments [post]
// @Security ApiKeyAuth
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	role, _ := c.Get("role")
	if roleStr, ok := role.(string); !ok || roleStr != "patient" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only patients can book appointments"})
		return
	}
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	patientID := uint(c.GetFloat64("user_id"))
	if existing, _ := h.Service.GetBySlotID(req.SlotID); existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slot already booked"})
		return
	}
	app := &model.Appointment{
		PatientID: patientID,
		DoctorID:  req.DoctorID,
		SlotID:    req.SlotID,
		Status:    "Upcoming",
	}
	if err := h.Service.Create(app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Book the slot in the schedule service
	go func(slotID uint) {
		url := "http://schedule:8082/slots/" + strconv.Itoa(int(slotID)) + "/book"
		req, _ := http.NewRequest("PUT", url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}(app.SlotID)
	c.JSON(http.StatusCreated, app)
}

func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	patientID := uint(c.GetFloat64("user_id"))
	apps, err := h.Service.GetByPatientID(patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment id"})
		return
	}
	app, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}
	var req struct {
		Status   string `json:"status"`
		SlotID   *uint  `json:"slot_id"`
		DoctorID *uint  `json:"doctor_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Handle rescheduling (slot/doctor change)
	if req.SlotID != nil && *req.SlotID != app.SlotID {
		// Unbook old slot
		go func(slotID uint) {
			url := "http://schedule:8082/slots/" + strconv.Itoa(int(slotID)) + "/unbook"
			req, _ := http.NewRequest("PUT", url, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
		}(app.SlotID)
		// Check if new slot is available
		if existing, _ := h.Service.GetBySlotID(*req.SlotID); existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New slot already booked"})
			return
		}
		// Book new slot
		go func(slotID uint) {
			url := "http://schedule:8082/slots/" + strconv.Itoa(int(slotID)) + "/book"
			req, _ := http.NewRequest("PUT", url, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
		}(*req.SlotID)
		app.SlotID = *req.SlotID
	}
	if req.DoctorID != nil {
		app.DoctorID = *req.DoctorID
	}
	if req.Status != "" {
		app.Status = req.Status
	}
	if err := h.Service.Update(app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, app)
}

func (h *AppointmentHandler) DeleteAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appointment id"})
		return
	}
	app, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
		return
	}
	if err := h.Service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Unbook the slot in the schedule service
	go func(slotID uint) {
		url := "http://schedule:8082/slots/" + strconv.Itoa(int(slotID)) + "/unbook"
		req, _ := http.NewRequest("PUT", url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// Optionally log error
			return
		}
		defer resp.Body.Close()
	}(app.SlotID)
	c.JSON(http.StatusOK, gin.H{"message": "appointment cancelled"})
}
