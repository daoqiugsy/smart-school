// service/ai_service.go (最终解析版)

package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smart-school/pkg/config"
	"strings"
)

// ... 结构体和接口定义不变 ...
type ChatStreamResponse struct {
	Content string
	Err     error
}
type AIService interface {
	ChatStream(ctx context.Context, query, userID, userType string) (<-chan ChatStreamResponse, error)
}
type aiService struct {
	cfg *config.CozeConfig
}

func NewAIService(cfg *config.CozeConfig) AIService {
	return &aiService{cfg: cfg}
}

// --- 定义用于解析 Coze 响应的结构体 ---
// CozeMessage 对应 data 字段的整体结构
type CozeMessage struct {
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	// ... 可以根据需要添加其他字段，如 NodeIsFinish
}

// CozeContent 对应 content 字段内部的 JSON 结构
type CozeContent struct {
	Output string `json:"output"`
}

func (s *aiService) ChatStream(ctx context.Context, query, userID, userType string) (<-chan ChatStreamResponse, error) {
	outChan := make(chan ChatStreamResponse, 1)

	go func() {
		defer close(outChan)

		// ... 请求构建部分完全不变 ...
		requestBody := map[string]interface{}{
			"parameters": map[string]interface{}{
				"query":     query,
				"user_id":   userID,
				"user_type": userType,
			},
			"workflow_id": s.cfg.WorkflowID,
		}
		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			outChan <- ChatStreamResponse{Err: fmt.Errorf("failed to marshal request body: %w", err)}
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.URL, bytes.NewBuffer(requestJSON))
		if err != nil {
			outChan <- ChatStreamResponse{Err: fmt.Errorf("failed to create request: %w", err)}
			return
		}

		req.Header.Set("Authorization", "Bearer "+s.cfg.Token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream")

		client := &http.Client{}
		resp, err := client.Do(req)
		// ... 错误处理部分完全不变 ...
		if err != nil { /* ... */
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK { /* ... */
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		scannerBuf := make([]byte, 64*1024)
		scanner.Buffer(scannerBuf, 64*1024)

		var currentEvent string

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "event:") {
				currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
				continue // 读取下一行
			}

			if strings.HasPrefix(line, "data:") {
				// 只处理我们关心的 Message 事件
				if currentEvent == "Message" {
					dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

					// 第一次解析：将 data 字符串解析到 CozeMessage 结构体
					var msg CozeMessage
					err := json.Unmarshal([]byte(dataStr), &msg)
					if err != nil {
						log.Printf("[WARN] Failed to unmarshal Coze message: %v. Raw data: %s", err, dataStr)
						continue // 解析失败，跳过这一条
					}

					// 检查 content_type 是否为 text
					if msg.ContentType == "text" {
						// 第二次解析：将 content 字符串解析到 CozeContent 结构体
						var content CozeContent
						err := json.Unmarshal([]byte(msg.Content), &content)
						if err != nil {
							log.Printf("[WARN] Failed to unmarshal inner content: %v. Raw content: %s", err, msg.Content)
							continue // 解析失败，跳过
						}

						// 成功提取最终的 AI 回复！
						finalOutput := content.Output
						if finalOutput != "" {
							// 为了让前端更容易处理，我们自己也包装成一个简单的 JSON
							responsePayload := gin.H{"text": finalOutput}
							responseJSON, _ := json.Marshal(responsePayload)

							outChan <- ChatStreamResponse{Content: string(responseJSON)}
						}
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			outChan <- ChatStreamResponse{Err: fmt.Errorf("error reading stream: %w", err)}
		}
	}()

	return outChan, nil
}
