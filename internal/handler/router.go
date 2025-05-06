package handler

import (
	"github.com/gin-gonic/gin"
	"smart-school/internal/middleware"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler) {
	// 公开路由组
	public := r.Group("/api")
	{
		// 认证相关路由
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	// 需要认证的路由组
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// 后续添加需要认证的路由
	}
}
