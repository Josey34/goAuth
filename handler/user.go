package handler

import (
	"goauth/dto"
	"goauth/errors"
	"goauth/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.UserUsecase.GetProfile(userID)
	if err != nil {
		statusCode := errors.ToHTTPStatus(err)
		c.JSON(statusCode, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userResp := dto.FromEntity(user)
	c.JSON(http.StatusOK, userResp)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		statusCode := errors.ToHTTPStatus(err)
		fieldErrors := extractFieldErrors(err)
		c.JSON(statusCode, dto.ErrorResponse{Error: "validation failed", Details: fieldErrors})
		return
	}

	userID := c.GetString("userID")
	user, err := h.UserUsecase.UpdateProfile(userID, req.Name)
	if err != nil {
		statusCode := errors.ToHTTPStatus(err)
		c.JSON(statusCode, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userResp := dto.FromEntity(user)
	c.JSON(http.StatusOK, userResp)
}

func (h *UserHandler) AdminDashboard(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to admin dashboard",
		"user_id": userID,
		"data":    "Admin-only content",
	})
}
