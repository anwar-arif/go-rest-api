package main

import (
	"context"
	infra "go-rest-api/infra/db"
	infraRedis "go-rest-api/infra/redis"
	"go-rest-api/logger"
	"log"

	"go-rest-api/config"
	infraMongo "go-rest-api/infra/mongo"
	"go-rest-api/repo"

	"github.com/spf13/cobra"
)

var repos []repo.Repo

var migrationRoot = &cobra.Command{
	Use:   "migration",
	Short: "Run database migrations",
	Long:  `Migration is a tool to generate and modify database tables`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfgMongo := config.GetMongo(cfgPath, env, envPath)
		cfgDBTable := config.GetTable(cfgPath)
		cfgRedis := config.GetRedis(cfgPath, envPath)

		ctx := context.Background()

		lgr := logger.DefaultOutStructLogger

		mgo, err := infraMongo.New(ctx, cfgMongo)
		if err != nil {
			return err
		}
		defer mgo.Close(ctx)

		rds, err := infraRedis.New(ctx, cfgRedis, lgr)
		if err != nil {
			return err
		}
		defer rds.Close()

		db := infra.NewDB(mgo, rds)

		brandRepo := repo.NewBrand(cfgDBTable.BrandCollectionName, db)
		userRepo := repo.NewUser(cfgDBTable.UserCollectionName, db)

		repos = []repo.Repo{
			brandRepo,
			userRepo,
		}

		return nil
	}}

func init() {
	migrationRoot.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yml", "config file path")
}

var migrationUp = &cobra.Command{
	Use:   "up",
	Short: "Populate tables in database",
	Long:  `Populate tables in database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Populating database indices...")
		for _, t := range repos {
			if err := t.EnsureIndices(); err != nil {
				log.Println(err)
			}
		}
		log.Println("Populating database indices successfully...")
		return nil
	},
}

var migrationDown = &cobra.Command{
	Use:   "down",
	Short: "Drop tables from database",
	Long:  `Drop tables from database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Dropping database table...")
		for _, t := range repos {
			if err := t.DropIndices(); err != nil {
				log.Println(err)
			}
		}

		log.Println("Database dropped successfully!")
		return nil
	},
}

func init() {
	migrationRoot.AddCommand(
		migrationUp,
		migrationDown,
	)
}
