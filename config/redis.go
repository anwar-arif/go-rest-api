package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

type Redis struct {
	URL       string        `json:"url"`
	DbID      int           `json:"db_id"`
	DBTimeOut time.Duration `json:"db_timeout"`
}

var redisConfig *Redis
var redisOnce sync.Once

func loadRedis(fileName string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	redisConfig = &Redis{
		URL:       viper.GetString("redis.url"),
		DbID:      viper.GetInt("redis.db_id"),
		DBTimeOut: viper.GetDuration("redis.time_out") * time.Second,
	}

	log.Println("redis config ", redisConfig)
	return nil
}

func GetRedis(fileName string) *Redis {
	redisOnce.Do(func() {
		loadRedis(fileName)
	})

	return redisConfig
}
