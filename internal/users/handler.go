package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaelassa01/gomacbot/internal/users/domain"
	u "github.com/michaelassa01/gomacbot/utils"
)

// Handler binds Gin routes to the service.
type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}



func (h *Handler) CreateUser(ctx *gin.Context) {
	// Implementation of CreateUser handler goes here
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := u.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:          req.Email,
		Username:       req.Username,
		Age:            req.Age,
		PhoneNumber:    req.PhoneNumber,
		Avatar:         req.Avatar,
		HashedPassword: hashedPassword,
	}
	createdUser, err := h.service.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, toUserResponse(createdUser))
}
