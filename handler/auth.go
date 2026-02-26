package handler

import (
	"goauth/dto"
	"goauth/errors"
	"goauth/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		statusCode := errors.ToHTTPStatus(err)
		c.JSON(statusCode, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.authUsecase.Register(req.Email, req.Password, req.Name)
	if err != nil {
		statusCode := errors.ToHTTPStatus(err)
		c.JSON(statusCode, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userResp := dto.FromEntity(user)

	c.JSON(http.StatusCreated, userResp)
}
