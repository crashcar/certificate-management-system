package ws

import (
	"fmt"
	"log"
	"net/http"

	"application/model"
	"application/pkg/app"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var clients = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade WebSocket: %v", err)
			appG.Response(http.StatusInternalServerError, "升级ws连接出错", err)
			return
		}
		defer ws.Close()

		// 将新的客户端连接添加到clients全局变量中
		userID := c.Query("userID")
		clients[userID] = ws

		// 发送未读消息给用户（用户刚登录时建立连接并升级连接，从数据库找未读消息发送）
		sendUnreadNotifications(userID, ws, db, appG)

		// for循环处理消息，直到一方断开连接跳出循环，然后跳出函数
		for {
			var message []byte
			// 尝试读取新消息
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Printf("read error: %v", err)
				// TODO: 加锁
				delete(clients, userID) // 在这里处理用户断开连接
				break
			}
			log.Printf("Received: %s", string(message))
		}
	}
}

func sendUnreadNotifications(userID string, conn *websocket.Conn, db *gorm.DB, appG app.Gin) {
	var notifications []model.Notification
	db.Where("user_id = ? AND is_read = false", userID).Find(&notifications)

	for _, notification := range notifications {
		err := conn.WriteJSON(notification)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "发送未读通知时出错", err)
			break
		}
	}
}

func WriteDBAndNotifyUser(db *gorm.DB, notification model.Notification, appG app.Gin) {
	// 创建通知，写入数据库
	result := db.Create(&notification)
	if result.Error != nil {
		appG.Response(http.StatusInternalServerError, "保存通知到数据库出错", result.Error.Error())
		return
	}

	// 检查用户是否在线
	userID := notification.UserID
	if ws, ok := clients[userID]; ok {
		// 用户在线，直接发送通知
		err := ws.WriteJSON(notification)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "在线用户发送通知时出错", err)
			return
		}
	}
}
