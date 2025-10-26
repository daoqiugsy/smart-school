package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"log"
	"os"
	"smart-school/internal/service"
	"smart-school/pkg/config"
)

type WSMessage struct {
	Query    string `json:"query"`
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
}

func main() {
	// Get config
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal("worker 加载配置文件失败：v%", err)
	}
	// 初始化AIServive
	aiService := service.NewAIService(&cfg.AI.Coze)
	// 连接RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("worker 连接RabbitMQ失败：v%", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("worker 获取Channel失败：v%", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ai_tasks",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("worker 声明Queue失败：v%", err)
	}

	// 从队列消费消息
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	// --- 新增：为 Worker 初始化 Redis 连接 ---
	redisAddr := "localhost:6379"
	if host := os.Getenv("REDIS_HOST"); host != "" {
		redisAddr = host + ":6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Worker 无法连接到 Redis: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("[✓]Received a message 接收到消息: %s", d.Body)

			var msg WSMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("无法解析消息：%s", err)
				continue
			}

			// 核心逻辑
			// 1. 调用AIService的ChatStream方法
			streamChan, err := aiService.ChatStream(context.Background(), msg.Query, msg.UserID, msg.UserType)
			if err != nil {
				log.Printf("调用AI服务失败：%s", err)
				continue
			}
			// 2. 从 channel 中读取结果，并通过WebSocket 推送用户
			for resp := range streamChan {
				var pushData []byte
				if resp.Err != nil {
					//// 错误也推送给用户
					// service.PushToWsClient(msg.UserID, websocket.TextMessage, []byte(resp.Err.Error()))
					//break
					pushData = []byte(resp.Err.Error())
					// ！！！关键改变：不再调用本地的 ws_manager ！！！
					// 将结果发布到 Redis 频道
					service.PublishToUser(rdb, msg.UserID, pushData)
					break
				} else {
					pushData = []byte(resp.Content)
					// ！！！关键改变：将结果发布到 Redis 频道 ！！！
					service.PublishToUser(rdb, msg.UserID, pushData)
				}
				//// 将AI的回复推送给指定用户
				// service.PushToWsClient(msg.UserID, websocket.TextMessage, []byte(resp.Content))
			}
			// 推送结束标志
			//service.PushToWsClient(msg.UserID, websocket.TextMessage, []byte("done"))
			service.PublishToUser(rdb, msg.UserID, []byte("[DONE]"))
			log.Println("[✓]Done Task for user %s finished", msg.UserID)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
