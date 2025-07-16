package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smart-school/internal/service"
	_ "strconv"
)

// ScheduleHandler 课程表处理器
type ScheduleHandler struct {
	scheduleService service.ScheduleService
}

// NewScheduleHandler 创建课程表处理器实例
func NewScheduleHandler(scheduleService service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

// ImportFromCSV 从CSV导入课程表
func (h *ScheduleHandler) ImportFromCSV(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传CSV文件"})
		return
	}
	defer file.Close()

	// 调用服务层导入课程表
	err = h.scheduleService.ImportFromCSV(uint(userID.(uint)), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程表导入成功"})
}

// ImportFromAPI 从教务系统API导入课程表
func (h *ScheduleHandler) ImportFromAPI(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取请求参数
	var req struct {
		APIURL   string `json:"api_url" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务层导入课程表
	err := h.scheduleService.ImportFromAPI(uint(userID.(uint)), req.APIURL, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传CSV文件"})
		return
	}
	defer file.Close()
}

// ImportFromExcel 从Excel导入课程表
func (h *ScheduleHandler) ImportFromExcel(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传Excel文件"})
		return
	}
	defer file.Close()

	// 调用服务层导入课程表
	err = h.scheduleService.ImportFromExcel(uint(userID.(uint)), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "课程表导入成功"})
}

// GetStudentSchedule 获取学生课程表
func (h *ScheduleHandler) GetStudentSchedule(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用服务层获取课程表
	schedules, err := h.scheduleService.GetStudentSchedule(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedules": schedules})
}
