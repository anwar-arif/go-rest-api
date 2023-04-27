package infra

import (
	"go-rest-api/infra/mongo"
	"go-rest-api/infra/redis"
)

type DB struct {
	Mongo *mongo.Mongo
	Redis *redis.Redis
}

func NewDB(mgo *mongo.Mongo, rds *redis.Redis) *DB {
	return &DB{
		Mongo: mgo,
		Redis: rds,
	}
}
