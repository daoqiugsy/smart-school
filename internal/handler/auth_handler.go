package handler

import (
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 调用服务层进行注册
	err := h.authService.Register(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := 500
		if err == service.ErrUserExists {
			statusCode = http.StatusConflict
			code = 409
		}
		c.JSON(statusCode, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": nil,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 调用服务层进行登录
	user, err := h.authService.Login(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		code := 500
		switch err {
		case service.ErrUserNotFound, service.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
			code = 401
		case service.ErrUserDisabled:
			statusCode = http.StatusForbidden
			code = 403
		}
		c.JSON(statusCode, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成令牌失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}
