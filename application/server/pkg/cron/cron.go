package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

const spec = "0 0 0 * * ?" // 每天0点执行
//const spec = "*/10 * * * * ?" //10秒执行一次，用于测试

func Init() {
	c := cron.New(cron.WithSeconds()) //支持到秒级别
	_, err := c.AddFunc(spec, GoRun)
	if err != nil {
		log.Printf("定时任务开启失败 %s", err)
	}
	c.Start()
	log.Printf("定时任务已开启")
	select {}
}

func GoRun() {
	log.Printf("定时任务已启动")
	// //先把所有销售查询出来
	// resp, err := bc.ChannelQuery("queryCertByInfosLists", [][]byte{}) //调用智能合约
	// if err != nil {
	// 	log.Printf("定时任务-queryCertByInfosLists失败%s", err.Error())
	// 	return
	// }
	// // 反序列化json
	// var data []model.LedgerCertificate
	// if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
	// 	log.Printf("定时任务-反序列化json失败%s", err.Error())
	// 	return
	// }
	// for _, v := range data {
	// 	//把无效和已过期证书筛选出来
	// 	if v.Status == model.CertStatusConstant()["expired"] ||
	// 		v.Status == model.CertStatusConstant()["invalid"] {
	// 	}
	// }
}
