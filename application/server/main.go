package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "application/api/v1"
	"application/blockchain"
	"application/pkg/cron"
	"application/routers"
)

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal

	blockchain.Init()
	go cron.Init()

	// 连接pg并创建users/userprofiles表格

	if err := v1.InitUserdb(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8000)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routers.InitRouter(),
	}
	log.Printf("[info] start http server listening %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("start http server failed %s", err)
	}
}
