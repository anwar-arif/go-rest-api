package repo

import (
	"context"
	"go-rest-api/infra/db"
	"go-rest-api/infra/mongo"
	"go-rest-api/model"
	"time"
)

type UserRepo interface {
	Repo
	Create(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.GetUserByEmailResponse, error)
	GetAuthUserByEmail(ctx context.Context, email string) (*model.AuthUserData, error)
	StoreToken(ctx context.Context, email string, token string) error
	GetToken(ctx context.Context, key string) (string, error)
	RemoveToken(ctx context.Context, key ...string) error
}

type UserDB struct {
	table string
	db    *infra.DB
}

func NewUser(table string, db *infra.DB) UserRepo {
	return &UserDB{
		table: table,
		db:    db,
	}
}

func (*UserDB) Indices() []mongo.DbIndex {
	return []mongo.DbIndex{
		{
			Name: "email",
			Keys: []mongo.DbIndexKey{
				{"email", 1},
			},
		},
	}
}

func (u *UserDB) EnsureIndices() error {
	return u.db.Mongo.EnsureIndices(context.Background(), u.table, u.Indices())
}

func (u *UserDB) DropIndices() error {
	return u.db.Mongo.DropIndices(context.Background(), u.table, u.Indices())
}

func (u *UserDB) Create(ctx context.Context, user *model.User) error {
	return u.db.Mongo.Insert(ctx, u.table, user)
}

func (u *UserDB) GetUserByEmail(ctx context.Context, email string) (*model.GetUserByEmailResponse, error) {
	q := mongo.DbQuery{
		{"email", email},
	}

	user := &model.GetUserByEmailResponse{}
	if err := u.db.Mongo.FindOne(ctx, u.table, q, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserDB) GetAuthUserByEmail(ctx context.Context, email string) (*model.AuthUserData, error) {
	q := mongo.DbQuery{
		{"email", email},
	}
	user := &model.AuthUserData{}

	if err := u.db.Mongo.FindOne(ctx, u.table, q, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserDB) StoreToken(ctx context.Context, email string, token string) error {
	return u.db.Redis.Set(ctx, email, token, time.Second*300)
}

func (u *UserDB) GetToken(ctx context.Context, key string) (string, error) {
	return u.db.Redis.Get(ctx, key)
}

func (u *UserDB) RemoveToken(ctx context.Context, key ...string) error {
	return u.db.Redis.Del(ctx, key...)
}
