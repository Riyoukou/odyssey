package repository

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	RedisClient *redis.Client
	DB          *gorm.DB
	err         error
)

func Init() {
	//initRedis()
	//initMysql()
	initSQLite()
}

/*func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	if err := RedisClient.Ping().Err(); err != nil {
		panic(err)
	}
	log.Println("Redis connected successfully to", viper.GetString("redis.address"))
}*/

/*
	func initMysql() {
		DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")),
			&gorm.Config{
				PrepareStmt:            true,
				SkipDefaultTransaction: true,
			},
		)
		if err != nil {
			panic(err)
		}
		log.Println("MySQL connected successfully to", viper.GetString("mysql.dsn"))
	}
*/
func initSQLite() {
	dbPath := viper.GetString("sqlite.path")

	if dbPath == "" {
		dbPath = "./data/odyssey.db"
	}

	dbDir := filepath.Dir(dbPath)

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		panic(err)
	}

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		panic(err)
	}

	log.Println("SQLite database:", absPath)

	DB, err = gorm.Open(
		sqlite.Open(absPath),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)

	if err != nil {
		panic(err)
	}

	log.Println("SQLite connected successfully:", absPath)

	// 自动创建数据库表
	if err := DB.AutoMigrate(
		&model.UserTable{},
	); err != nil {
		panic(err)
	}

	initDefaultUsers()

	log.Println("SQLite tables migrated successfully")
}

func initDefaultUsers() {
	var count int64

	if err := DB.Model(&model.UserTable{}).Count(&count).Error; err != nil {
		panic(err)
	}

	// 已经存在用户，不再初始化
	if count > 0 {
		return
	}

	users := []model.UserTable{
		{
			ID:       1,
			Name:     "admin",
			Email:    "admin@odyssey.cn",
			Phone:    "00000000000",
			Token:    "",
			Password: "$2a$10$6zsDHgyUq2/098MNsiwMw.dlVYMeWUMjyZGgCxOnfNpJ28ANDgZsC",
			Type:     "local",
			Role:     "admin",
		},
		{
			ID:       2,
			Name:     "user",
			Email:    "user@odyssey.cn",
			Phone:    "00000000000",
			Token:    "",
			Password: "$2a$10$6zsDHgyUq2/098MNsiwMw.dlVYMeWUMjyZGgCxOnfNpJ28ANDgZsC",
			Type:     "local",
			Role:     "user",
		},
	}

	if err := DB.Create(&users).Error; err != nil {
		panic(err)
	}

	log.Println("Default users created successfully")
}
