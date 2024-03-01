package v1

import (
	"application/model"
	"application/pkg/app"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReadNotificationRequestBody struct {
	NotificationID uint `json:"notificationID"`
}

// 用户已读通知标记
func ReadNotification(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		body := new(ReadNotificationRequestBody)
		//解析Body参数
		if err := c.ShouldBind(body); err != nil {
			appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
			return
		}
		if body.NotificationID == 0 {
			appG.Response(http.StatusBadRequest, "失败", "notificationID不能为0")
			return
		}

		var notification model.Notification
		// 标记对应通知为已读
		result := db.Model(&notification).Where("id = ?", body.NotificationID).Update("is_read", true)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库已读标记出错", result.Error.Error())
			return
		}

		appG.Response(http.StatusOK, "成功", "数据库已读标记成功")
	}
}
