package db

import (
	"time"

	"github.com/xuanbo/ohmydata/pkg/config"
	orm "github.com/xuanbo/ohmydata/pkg/db/gorm"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB gorm db
	DB   *gorm.DB
	node *snowflake.Node
)

// Init 初始化数据库
func Init() error {
	url := config.GetString("mysql.url")
	maxIdleConns := config.GetInt("mysql.maxIdleConns")
	maxOpenConns := config.GetInt("mysql.maxOpenConns")

	log.Logger().Info("初始化数据库", zap.String("url", url))

	var err error
	DB, err = gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: orm.NewZapLogger(log.Logger(), 200*time.Millisecond),
	})
	if err != nil {
		return err
	}
	// 设置日志级别
	DB.Logger.LogMode(logger.Info)
	// 设置连接池
	db, err := DB.DB()
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err := db.Ping(); err != nil {
		return err
	}

	// Create a new Node with a Node number of 1
	node, err = snowflake.NewNode(1)
	if err != nil {
		return err
	}

	// 同步表结构
	if err := syncDB(); err != nil {
		return err
	}

	return nil
}

func syncDB() error {
	if err := DB.AutoMigrate(
		&entity.User{},
		&entity.DataSource{}, &entity.DataSet{},
		&entity.RequestParam{}, &entity.ResponseParam{},
	); err != nil {
		return err
	}

	var user entity.User
	if err := DB.Where("username = ?", "admin").Find(&user).Error; err != nil {
		return err
	}
	if user.ID == "" {
		user.ID = NewID()
		user.Name = "管理员"
		user.Username = "admin"
		user.Password = "123456"
		user.CreatedBy = user.ID
		if err := DB.Create(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

// NewID Generate a snowflake ID.
func NewID() string {
	id := node.Generate()
	return id.String()
}
