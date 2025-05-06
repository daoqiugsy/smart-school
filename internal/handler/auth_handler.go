package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smart-school/internal/service"
	"smart-school/pkg/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Register")
	// 调用服务层进行注册
	err := h.authService.Register(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrUserExists {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层进行登录
	user, err := h.authService.Login(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case service.ErrUserNotFound, service.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
		case service.ErrUserDisabled:
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
