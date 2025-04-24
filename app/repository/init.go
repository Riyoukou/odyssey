package repository

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	RedisClient *redis.Client
	DB          *gorm.DB
	err         error
)

func Init() {
	initRedis()
	initMysql()
}

func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	if err := RedisClient.Ping().Err(); err != nil {
		panic(err)
	}
	log.Println("Redis connected successfully to", viper.GetString("redis.address"))
}

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
