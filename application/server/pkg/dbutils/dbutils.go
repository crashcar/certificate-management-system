package dbutils

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"application/model"
)

const (
	maxRetries    = 5                // 最大重试次数
	retryInterval = 10 * time.Second // 重试间隔
)

// const createUsersTableSQL = `CREATE TABLE IF NOT EXISTS users (
//     id SERIAL PRIMARY KEY,
//     real_name VARCHAR(255) NOT NULL,
//     password VARCHAR(255) NOT NULL,
//     email VARCHAR(255) NOT NULL,
//     role INTEGER DEFAULT 0
// );`

// const createCertsTableSQL = `CREATE TABLE IF NOT EXISTS certs (
//     id SERIAL PRIMARY KEY,
//     url VARCHAR(255) NOT NULL,
// 	cert_type VARCHAR(255) NOT NULL,
//     created_at DATETIME NOT NULL,
//     uploader_id VARCHAR(255) NOT NULL,
// 	uploader_name VARCHAR(255) NOT NULL,
//     is_processed BOOLEAN DEFAULT FALSE
// );`

const createSeqSQL = "CREATE SEQUENCE IF NOT EXISTS admin_id_seq START WITH 10001;"

const triggerFuncSQL = `CREATE OR REPLACE FUNCTION increment_admin_id()
RETURNS TRIGGER AS $$
BEGIN
    NEW.id := nextval('admin_id_seq');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;`

const createTriggerSQL = `CREATE TRIGGER set_admin_id_before_insert
BEFORE INSERT ON admins
FOR EACH ROW
EXECUTE FUNCTION increment_admin_id();`

func NewDB() (*gorm.DB, error) {
	//todo: 数据库连接命令
	dsn := "host=certman-db user=admin password=1234 dbname=userdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var db *gorm.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("Failed to connect to database. Retry %d times: %v", i+1, err)
		time.Sleep(retryInterval)
	}

	return nil, err
}

// 建表：users-用户表 tempCerts-待处理的证书暂存表
func InitDB(db *gorm.DB) error {
	// 手动建表
	// if err := db.Exec(createUsersTableSQL).Error; err != nil {
	// 	return err
	// }

	// if err := db.Exec(createCertsTableSQL).Error; err != nil {
	// 	return err
	// }

	// 自动迁移模式建表
	err := db.AutoMigrate(&model.User{}, &model.Cert{}, &model.Admin{})
	if err != nil {
		return err
	}

	// 设置 admins 表格的主键开始序号为 10001
	if err = db.Exec(createSeqSQL).Error; err != nil {
		return err
	}
	if err = db.Exec(triggerFuncSQL).Error; err != nil {
		return err
	}
	// 检查触发器是否存在
	var triggerCount int64
	if err := db.Raw("SELECT count(*) FROM pg_trigger WHERE tgname = 'set_admin_id_before_insert'").Count(&triggerCount).Error; err != nil {
		return err
	}

	// 如果触发器不存在，则创建
	if triggerCount == 0 {
		if err = db.Exec(createTriggerSQL).Error; err != nil {
			return err
		}
	}
	return nil
}
