// service/redis_pubsub.go
package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
)

// SubscribeAndForward 订阅指定用户的 Redis 频道，并将消息转发到 WebSocket
func SubscribeAndForward(ctx context.Context, userID string, ws *websocket.Conn) {
	// 为这个函数也创建一个独立的 Redis 连接
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"}) // 简单起见，地址硬编码
	defer rdb.Close()

	channelName := fmt.Sprintf("ws:user:%s", userID)
	pubsub := rdb.Subscribe(ctx, channelName)
	defer pubsub.Close()

	log.Printf("User %s subscribed to Redis channel '%s'", userID, channelName)

	// 等待订阅确认
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Printf("Redis subscribe error for user %s: %v", userID, err)
		return
	}

	// 创建一个 channel 来接收 Redis 消息
	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done(): // 如果 Gin 的请求上下文结束（例如客户端断开连接）
			log.Printf("Context done for user %s, unsubscribing.", userID)
			return
		case msg := <-ch: // 从 Redis 频道收到了消息
			// 将收到的消息写入 WebSocket
			err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				log.Printf("Failed to write message to user %s: %v", userID, err)
				return // 写入失败，结束这个 goroutine
			}
		}
	}
}

// PublishToUser 将消息发布到指定用户的 Redis 频道
func PublishToUser(rdb *redis.Client, userID string, data []byte) {
	channelName := fmt.Sprintf("ws:user:%s", userID)
	err := rdb.Publish(context.Background(), channelName, data).Err()
	if err != nil {
		log.Printf("Failed to publish message to user %s on channel %s: %v", userID, channelName, err)
	}
}
