package database

import (
	"competitive-analyzer/models"
	"log"

	"github.com/glebarez/sqlite" // 纯Go SQLite驱动，无需CGO
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&models.DiscoveryTask{},
		&models.SearchCache{},
		&models.Competitor{},
		&models.DataSource{},
		&models.RawContent{},
		&models.ParsedData{},
		&models.AnalysisReport{},
		&models.ChangeLog{},
		&models.MonitorTask{},
	)
	if err != nil {
		return err
	}

	log.Println("数据库初始化成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
