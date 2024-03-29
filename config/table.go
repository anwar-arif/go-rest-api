package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// Table holds table configurations
type Table struct {
	BrandCollectionName string `yaml:"brand"`
	UserCollectionName  string `yaml:"user"`
}

var tableOnce = sync.Once{}
var tableConfig *Table

// loadTable loads config from path
func loadTable(fileName string) error {
	readConfig(fileName)
	viper.AutomaticEnv()

	tableConfig = &Table{
		BrandCollectionName: viper.GetString("collection.brand"),
		UserCollectionName:  viper.GetString("collection.user"),
	}

	log.Println("table config ", tableConfig)

	return nil
}

// GetTable returns table config
func GetTable(fileName string) *Table {
	tableOnce.Do(func() {
		loadTable(fileName)
	})

	return tableConfig
}
