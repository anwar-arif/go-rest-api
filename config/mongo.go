package config

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Mongo holds mongo config
type Mongo struct {
	Host      string        `yaml:"host"`
	Port      string        `yaml:"port"`
	UserName  string        `yaml:"user_name"`
	Password  string        `yaml:"password"`
	URL       string        `yaml:"url"`
	DBName    string        `yaml:"db_name"`
	DBTimeOut time.Duration `yaml:"time_out"`
}

func (m *Mongo) url(env string) string {
	if env == "dev" || env == "test" {
		return "mongodb://" + m.UserName + ":" + m.Password + "@" + m.Host + ":" + m.Port
	}
	return "mongodb://" + m.Host + ":" + m.Port
}

var mongoOnce = sync.Once{}
var mongoConfig *Mongo

// loadMongo loads config from path
func loadMongo(fileName, env, envFilePath string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	hostEnv := viper.GetString("mongo.env.host")
	portEnv := viper.GetString("mongo.env.port")
	dbUserNameEnv := viper.GetString("mongo.env.user_name")
	dbPasswordEnv := viper.GetString("mongo.env.password")
	dbName := viper.GetString("mongo.db_name")
	dbTimeOut := viper.GetDuration("mongo.time_out") * time.Second

	readConfig(envFilePath)
	viper.AutomaticEnv()

	mongoConfig = &Mongo{
		Host:      viper.GetString(hostEnv),
		Port:      viper.GetString(portEnv),
		UserName:  viper.GetString(dbUserNameEnv),
		Password:  viper.GetString(dbPasswordEnv),
		DBName:    dbName,
		DBTimeOut: dbTimeOut,
	}

	mongoConfig.URL = mongoConfig.url(env)

	log.Println("dbUserNameEnv ", dbUserNameEnv)
	log.Println("env", env)
	log.Println("envFile", envFilePath)
	log.Println("mongo config ", mongoConfig)
	return nil
}

// GetMongo returns mongo config
func GetMongo(fileName, env, envFilePath string) *Mongo {
	mongoOnce.Do(func() {
		loadMongo(fileName, env, envFilePath)
	})

	return mongoConfig
}
