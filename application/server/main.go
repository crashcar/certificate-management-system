package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"application/blockchain"
	"application/pkg/cron"
	"application/pkg/dbutils"
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

	// 连接pg并创建表格
	db, err := dbutils.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to pg:", err)
	}

	if err := dbutils.InitDB(db); err != nil {
		log.Fatal("Failed to initiate database:", err)
	}

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8000)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routers.InitRouter(db),
	}
	log.Printf("[info] start http server listening %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("start http server failed %s", err)
	}
}
