package handler

import (
	"net/http"
	"smart-school/internal/ai"

	"github.com/gin-gonic/gin"
)

// AIHandler AI功能处理器
type AIHandler struct {
	courseAssistant *ai.CourseAssistant
}

// NewAIHandler 创建AI功能处理器
func NewAIHandler(courseAssistant *ai.CourseAssistant) *AIHandler {
	return &AIHandler{
		courseAssistant: courseAssistant,
	}
}

// CourseAssistantQuery 课程管理小助手查询请求
type CourseAssistantQuery struct {
	Query string `json:"query" binding:"required"`
}

// HandleCourseAssistantQuery 处理课程管理小助手查询
func (h *AIHandler) HandleCourseAssistantQuery(c *gin.Context) {
	// 从JWT获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 解析请求
	var req CourseAssistantQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理查询
	response, err := h.courseAssistant.HandleQuery(c.Request.Context(), userID.(uint), req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
