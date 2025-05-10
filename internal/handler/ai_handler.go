package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"smart-school/pkg/config"

	"github.com/gin-gonic/gin"
)

type AIService interface {
	Chat(query string, studentID string) (string, error)
}

type AIHandler struct {
	aiService AIService
	config    *config.CozeConfig
}

func NewAIHandler(cfg *config.CozeConfig) *AIHandler {
	return &AIHandler{config: cfg}
}

// ChatRequest 定义聊天请求结构
type ChatRequest struct {
	Query string `json:"query" binding:"required"`
}

// Chat 处理AI聊天请求
func (h *AIHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的请求参数", "data": nil})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未找到用户信息", "data": nil})
		return
	}

	// 将 userID 转换为字符串
	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 创建一个通道来接收错误
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	// 在新的 goroutine 中调用 Coze API
	go func() {
		// 使用配置文件中的参数
		url := h.config.URL
		token := h.config.Token
		workflowID := h.config.WorkflowID

		// 构建请求体
		requestBody := map[string]interface{}{
			"parameters": map[string]interface{}{
				"query":      req.Query,
				"student_id": userIDStr,
			},
			"workflow_id": workflowID,
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			errChan <- err
			return
		}

		// 创建HTTP请求
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
		if err != nil {
			errChan <- err
			return
		}

		// 设置请求头
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		// 读取并解析响应
		scanner := bufio.NewScanner(resp.Body)
		var event string

		for scanner.Scan() {
			line := scanner.Text()

			// 跳过空行
			if line == "" {
				continue
			}

			// 解析事件类型
			if strings.HasPrefix(line, "event: ") {
				event = strings.TrimPrefix(line, "event: ")
				continue
			}

			// 解析数据
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")

				// 只处理 Message 事件
				if event == "Message" {
					var sseData struct {
						Content      string `json:"content"`
						NodeIsFinish bool   `json:"node_is_finish"`
					}
					if err := json.Unmarshal([]byte(data), &sseData); err == nil {
						// 构建标准格式的响应
						response := gin.H{
							"code": 200,
							"msg":  "success",
							"data": sseData.Content,
						}

						// 将响应转换为 SSE 格式
						responseJSON, _ := json.Marshal(response)
						c.Writer.WriteString("data: " + string(responseJSON) + "\n\n")
						c.Writer.Flush()

						// 如果是最后一条消息，结束处理
						if sseData.NodeIsFinish {
							doneChan <- true
							return
						}
					}
				} else if event == "Done" {
					doneChan <- true
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
			return
		}
		doneChan <- true
	}()

	// 等待错误、完成或连接关闭
	select {
	case err := <-errChan:
		// 构建错误响应
		response := gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		}
		responseJSON, _ := json.Marshal(response)
		c.Writer.WriteString("data: " + string(responseJSON) + "\n\n")
		c.Writer.Flush()
	case <-doneChan:
		return
	case <-c.Request.Context().Done():
		return
	}
}

// callCozeAPI 调用Coze API
func (h *AIHandler) callCozeAPI(query string, studentID string) (string, error) {
	// 使用配置文件中的参数
	url := h.config.URL
	token := h.config.Token
	workflowID := h.config.WorkflowID

	// 构建请求体
	requestBody := map[string]interface{}{
		"parameters": map[string]interface{}{
			"query":      query,
			"student_id": studentID,
		},
		"workflow_id": workflowID,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 这里简单返回响应内容，实际应用中可能需要解析JSON
	return string(respBody), nil
}
