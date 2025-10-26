package service

//
//import (
//	"github.com/gorilla/websocket"
//	"log"
//	"sync"
//)
//
//// 使用一个 map 存储客户端连接，key是userID
//var clients = make(map[string]*websocket.Conn)
//var mu sync.Mutex // 互斥锁并发安全
//
//// AddClient 向客户端连接添加一个客户端w
//func AddWsClient(userID string, conn *websocket.Conn) {
//	mu.Lock()
//	defer mu.Unlock()
//	clients[userID] = conn
//	log.Printf("WebSocket client added: %s", userID)
//}
//
//// RemoveClient 从客户端连接中删除一个客户端w
//func RemoveWsClient(userID string) {
//	mu.Lock()
//	defer mu.Unlock()
//	delete(clients, userID)
//	log.Printf("WebSocket client removed: [%s]", userID)
//}
//
//// PushToWsClient 向指定客户端推送消息
//func PushToWsClient(userID string, messageType int, data []byte) {
//	mu.Lock()
//	defer mu.Unlock()
//	log.Printf("[WS Manager] Attempting to push to userID: '%s'", userID)
//	if conn, ok := clients[userID]; ok {
//		log.Printf("[WS Manager] Found connection for userID: '%s'. Pushing data.", userID)
//		if err := conn.WriteMessage(messageType, data); err != nil {
//			log.Printf(" 写入 WebSocket 失败 for user [%s] %v", userID, err)
//			delete(clients, userID)
//		}
//	} else {
//		// --- 没找到连接！---
//		log.Printf("[WS Manager] No connection found for userID: '%s'. Message dropped.", userID)
//	}
//}
