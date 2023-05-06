package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

type Redis struct {
	Host      string        `yaml:"host"`
	Port      string        `yaml:"port"`
	URL       string        `yaml:"url"`
	DbID      int           `yaml:"db_id"`
	DBTimeOut time.Duration `yaml:"db_timeout"`
}

func (r *Redis) url() string {
	return r.Host + ":" + r.Port
}

var redisConfig *Redis
var redisOnce sync.Once

func loadRedis(fileName, envFilePath string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	redisHostEnv := viper.GetString("redis.env.host")
	redisPortEnv := viper.GetString("redis.env.port")
	dbId := viper.GetInt("redis.db_id")
	dbTimeOut := viper.GetDuration("redis.time_out") * time.Second

	readConfig(envFilePath)
	viper.AutomaticEnv()

	redisConfig = &Redis{
		Host:      viper.GetString(redisHostEnv),
		Port:      viper.GetString(redisPortEnv),
		DbID:      dbId,
		DBTimeOut: dbTimeOut,
	}
	redisConfig.URL = redisConfig.url()

	log.Println("redis config ", redisConfig)
	return nil
}

func GetRedis(fileName, envFilePath string) *Redis {
	redisOnce.Do(func() {
		loadRedis(fileName, envFilePath)
	})

	return redisConfig
}
