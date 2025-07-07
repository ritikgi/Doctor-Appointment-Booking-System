package handler

import (
	"net/http"
	"schedule/pkg/model"
	"schedule/pkg/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SlotHandler struct {
	Service *service.SlotService
}

// CreateSlot godoc
// @Summary Create a new slot
// @Description Doctor creates a new slot
// @Tags slots
// @Accept json
// @Produce json
// @Param slot body SlotRequest true "Slot info"
// @Success 201 {object} model.Slot
// @Failure 400 {object} map[string]string
// @Router /slots [post]
// @Security ApiKeyAuth
func (h *SlotHandler) CreateSlot(c *gin.Context) {
	role, _ := c.Get("role")
	if roleStr, ok := role.(string); !ok || roleStr != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only doctors can create slots"})
		return
	}
	var req struct {
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format"})
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format"})
		return
	}
	doctorID := uint(c.GetFloat64("user_id"))
	slot := &model.Slot{
		DoctorID:  doctorID,
		StartTime: start,
		EndTime:   end,
		IsBooked:  false, // New slots are available by default
	}
	if err := h.Service.Create(slot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, slot)
}

// GET /slots?doctor_id=... - get available slots for a doctor
func (h *SlotHandler) GetSlots(c *gin.Context) {
	doctorIDStr := c.Query("doctor_id")
	if doctorIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "doctor_id is required"})
		return
	}
	doctorID, err := strconv.Atoi(doctorIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id"})
		return
	}
	var slots []model.Slot
	if err := h.Service.Repo.DB.Where("doctor_id = ?", doctorID).Find(&slots).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch slots"})
		return
	}
	c.JSON(200, slots)
}

// PUT /slots/:id - doctor updates a slot
func (h *SlotHandler) UpdateSlot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}
	slot, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "slot not found"})
		return
	}
	var req struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.StartTime != "" {
		start, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format"})
			return
		}
		slot.StartTime = start
	}
	if req.EndTime != "" {
		end, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format"})
			return
		}
		slot.EndTime = end
	}
	if err := h.Service.Update(slot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slot)
}

// DELETE /slots/:id - doctor deletes a slot
func (h *SlotHandler) DeleteSlot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}
	if err := h.Service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "slot deleted"})
}

// PUT /slots/:id/book - mark a slot as booked
func (h *SlotHandler) BookSlot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid slot id"})
		return
	}
	var slot model.Slot
	if err := h.Service.Repo.DB.First(&slot, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Slot not found"})
		return
	}
	if slot.IsBooked {
		c.JSON(400, gin.H{"error": "Slot already booked"})
		return
	}
	slot.IsBooked = true
	if err := h.Service.Repo.DB.Save(&slot).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to book slot"})
		return
	}
	c.JSON(200, slot)
}

// PUT /slots/:id/unbook - mark a slot as available
func (h *SlotHandler) UnbookSlot(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid slot id"})
		return
	}
	var slot model.Slot
	if err := h.Service.Repo.DB.First(&slot, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Slot not found"})
		return
	}
	slot.IsBooked = false
	if err := h.Service.Repo.DB.Save(&slot).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to unbook slot"})
		return
	}
	c.JSON(200, slot)
}

// SlotRequest represents the request body for creating a slot
// swagger:model
type SlotRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
