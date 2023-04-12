package repo

import (
	"context"
	"go-rest-api/infra"
	"go-rest-api/model"
)

type UserRepo interface {
	Repo
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email *string) (*model.User, error)
}

type MgoUser struct {
	table string
	db    infra.DB
}

func NewUserRepo(table string, db infra.DB) UserRepo {
	return &MgoUser{
		table: table,
		db:    db,
	}
}

func (*MgoUser) Indices() []infra.DbIndex {
	return []infra.DbIndex{
		{
			Name: "email",
			Keys: []infra.DbIndexKey{
				{"email", 1},
			},
		},
	}
}

func (u *MgoUser) EnsureIndices() error {
	return u.db.EnsureIndices(context.Background(), u.table, u.Indices())
}

func (u *MgoUser) DropIndices() error {
	return u.db.DropIndices(context.Background(), u.table, u.Indices())
}

func (u *MgoUser) Create(ctx context.Context, user *model.User) error {
	return u.db.Insert(ctx, u.table, user)
}

func (u *MgoUser) GetByEmail(ctx context.Context, email *string) (*model.User, error) {
	q := infra.DbQuery{
		{"email", email},
	}
	user := &model.User{}

	if err := u.db.FindOne(ctx, u.table, q, &user); err != nil {
		return nil, err
	}
	return user, nil
}
