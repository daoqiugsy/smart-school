package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	// 设置响应头，确保流式传输正常工作
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache, no-transform")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx的缓冲
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// 清除所有可能的缓冲
	c.Writer.Flush()

	// 发送初始响应，确保连接已建立
	initialResponse := gin.H{
		"code": 200,
		"msg":  "connected",
		"data": "连接已建立，等待响应...",
	}
	initialJSON, _ := json.Marshal(initialResponse)
	c.Writer.WriteString("data: " + string(initialJSON) + "\n\n")
	c.Writer.Flush()

	// 等待一小段时间以确保初始响应能够送达
	time.Sleep(100 * time.Millisecond)

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
		req.Header.Set("Accept", "text/event-stream") // 明确请求流式响应

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		// 如果请求不成功，返回错误
		if resp.StatusCode != http.StatusOK {
			// 读取错误响应
			errBody, _ := io.ReadAll(resp.Body)
			errChan <- fmt.Errorf("API请求失败，状态码：%d, 响应: %s", resp.StatusCode, string(errBody))
			return
		}

		// 读取并解析响应
		scanner := bufio.NewScanner(resp.Body)
		// 增加scanner的缓冲区大小，避免长行问题
		scannerBuf := make([]byte, 64*1024)
		scanner.Buffer(scannerBuf, 64*1024)

		var event string
		var messageSent bool
		var messageCount int = 0

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
						messageSent = true
						messageCount++

						// 将消息拆分成更小的块，以确保流式效果
						// 对较长消息进行字符级拆分
						content := sseData.Content
						chunks := splitContent(content)

						for _, chunk := range chunks {
							if chunk == "" {
								continue
							}

							// 构建标准格式的响应
							response := gin.H{
								"code": 200,
								"msg":  "success",
								"data": chunk,
							}

							// 将响应转换为 SSE 格式
							responseJSON, _ := json.Marshal(response)
							c.Writer.WriteString("data: " + string(responseJSON) + "\n\n")
							c.Writer.Flush() // 关键：立即刷新缓冲区

							// 短暂延迟帮助客户端接收到分离的事件
							time.Sleep(10 * time.Millisecond)
						}

						// 如果是最后一条消息，结束处理
						if sseData.NodeIsFinish {
							doneChan <- true
							return
						}
					} else {
						// 解析JSON失败，发送原始数据
						fmt.Printf("解析JSON失败: %v, 原始数据: %s\n", err, data)
						c.Writer.WriteString("data: " + data + "\n\n")
						c.Writer.Flush()
					}
				} else if event == "Done" {
					doneChan <- true
					return
				} else {
					// 其他事件类型的处理
					fmt.Printf("收到其他事件: %s, 数据: %s\n", event, data)
					if !messageSent {
						// 如果尚未发送过消息，把原始数据发给客户端
						c.Writer.WriteString("data: " + data + "\n\n")
						c.Writer.Flush()
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
			return
		}

		// 如果没有发送过任何消息，可能是扫描器无法正确处理响应
		if !messageSent {
			// 尝试直接读取整个响应
			resp.Body.Close() // 关闭当前的body

			// 重新发起请求获取完整响应
			fullReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
			fullReq.Header.Set("Authorization", "Bearer "+token)
			fullReq.Header.Set("Content-Type", "application/json")

			fullResp, err := client.Do(fullReq)
			if err != nil {
				errChan <- err
				return
			}
			defer fullResp.Body.Close()

			respBody, _ := io.ReadAll(fullResp.Body)

			// 发送完整响应
			response := gin.H{
				"code": 200,
				"msg":  "full_response",
				"data": string(respBody),
			}
			responseJSON, _ := json.Marshal(response)
			c.Writer.WriteString("data: " + string(responseJSON) + "\n\n")
			c.Writer.Flush()
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
		// 发送完成信号
		completeResponse := gin.H{
			"code": 200,
			"msg":  "completed",
			"data": "",
		}
		completeJSON, _ := json.Marshal(completeResponse)
		c.Writer.WriteString("data: " + string(completeJSON) + "\n\n")
		c.Writer.Flush()
		return
	case <-c.Request.Context().Done():
		return
	}
}

// 将文本内容拆分为更小的块，以确保流式效果
func splitContent(content string) []string {
	if len(content) <= 10 {
		return []string{content}
	}

	// 尝试按标点符号拆分
	var chunks []string
	runes := []rune(content)

	// 定义分隔符，优先在这些位置断句
	delimiters := []rune{',', '，', '.', '。', '!', '！', '?', '？', ';', '；', ':', '：', '\n'}

	start := 0
	for i := 0; i < len(runes); i++ {
		// 检查当前字符是否是分隔符
		isDelimiter := false
		for _, d := range delimiters {
			if runes[i] == d {
				isDelimiter = true
				break
			}
		}

		// 分隔符处理：当是分隔符且当前子串长度>0时，或者当前子串已经足够长时
		if (isDelimiter && i-start > 0) || i-start >= 15 {
			if isDelimiter {
				// 包含分隔符
				chunks = append(chunks, string(runes[start:i+1]))
				start = i + 1
			} else {
				// 没有遇到分隔符但长度已经够长，尝试向前查找最近的分隔符
				foundBackDelimiter := false
				for j := i; j > start && j > i-5; j-- {
					for _, d := range delimiters {
						if runes[j] == d {
							chunks = append(chunks, string(runes[start:j+1]))
							start = j + 1
							i = j // 重新从分隔符后开始
							foundBackDelimiter = true
							break
						}
					}
					if foundBackDelimiter {
						break
					}
				}

				// 如果向前没找到合适的分隔符，就在当前位置切分
				if !foundBackDelimiter {
					chunks = append(chunks, string(runes[start:i+1]))
					start = i + 1
				}
			}
		}
	}

	// 处理剩余部分
	if start < len(runes) {
		chunks = append(chunks, string(runes[start:]))
	}

	// 如果没有成功拆分（可能没有合适的分隔符），则按固定长度拆分
	if len(chunks) == 0 {
		chunkSize := 15
		for i := 0; i < len(runes); i += chunkSize {
			end := i + chunkSize
			if end > len(runes) {
				end = len(runes)
			}
			chunks = append(chunks, string(runes[i:end]))
		}
	}

	return chunks
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
