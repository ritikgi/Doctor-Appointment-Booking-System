package handler

import (
	"net/http"
	"user/pkg/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

// GetProfile godoc
// @Summary Get current user's profile
// @Description Get current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /me [get]
// @Security ApiKeyAuth
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetFloat64("user_id") // JWT claims are float64
	user, err := h.Service.GetProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UserRequest represents the request body for updating a user profile
// swagger:model
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the current user
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRequest true "User info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetFloat64("user_id")
	user, err := h.Service.GetProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Name = req.Name
	if err := h.Service.UpdateProfile(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GET /doctors - list all doctors
func (h *UserHandler) ListDoctors(c *gin.Context) {
	doctors, err := h.Service.ListByRole("doctor")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}
