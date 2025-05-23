package handler

import (
	"smart-school/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, authHandler *AuthHandler, scheduleHandler *ScheduleHandler, aiHandler *AIHandler) {
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
		// 课程表相关路由
		schedule := protected.Group("/schedule")
		{
			// 导入课程表
			schedule.POST("/import/csv", scheduleHandler.ImportFromCSV)
			schedule.POST("/import/api", scheduleHandler.ImportFromAPI)
			// 获取学生课程表
			schedule.GET("/student", scheduleHandler.GetStudentSchedule)
			//ai助手
			schedule.POST("/ai/chat", aiHandler.Chat)
		}
	}
}
