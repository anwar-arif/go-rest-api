package repo

import (
	"context"
	"fmt"
	"go-rest-api/infra/db"
	"go-rest-api/infra/mongo"

	"go-rest-api/model"
)

// BrandRepo returns brand repo
type BrandRepo interface {
	Repo
	Create(ctx context.Context, bi *model.BrandInfo) error
	ListBrands(ctx context.Context, search string, skip, limit int64) ([]model.BrandInfo, error)
	GetBrandDetails(ctx context.Context, slug string) (*model.BrandInfo, error)
	GetBySlug(ctx context.Context, slug string) (*model.BrandInfo, error)
}

// MgoBrand brand repo
type MgoBrand struct {
	table string
	db    *infra.DB
}

// NewBrand returns new brand repo
func NewBrand(table string, db *infra.DB) BrandRepo {
	return &MgoBrand{
		table: table,
		db:    db,
	}
}

// Indices returns indices
func (*MgoBrand) Indices() []mongo.DbIndex {
	return []mongo.DbIndex{
		{
			Name: "slug_1_version_1",
			Keys: []mongo.DbIndexKey{
				{"slug", 1},
				{"version", 1},
			},
		},
		{
			Keys: []mongo.DbIndexKey{
				{"status", 1},
				{"approved", 1},
				{"categories", 1},
				{"name", 1},
				{"id", -1},
				{"brand_score", -1},
			},
		},
	}
}

// EnsureIndices ...
func (p *MgoBrand) EnsureIndices() error {
	return p.db.Mongo.EnsureIndices(context.Background(), p.table, p.Indices())
}

// DropIndices ...
func (p *MgoBrand) DropIndices() error {
	return p.db.Mongo.DropIndices(context.Background(), p.table, p.Indices())
}

// Create ...
func (p *MgoBrand) Create(ctx context.Context, bi *model.BrandInfo) error {
	return p.db.Mongo.Insert(ctx, p.table, bi)
}

// ListBrands ...
func (p *MgoBrand) ListBrands(ctx context.Context, search string, skip, limit int64) ([]model.BrandInfo, error) {
	query := mongo.DbQuery{
		{"status", model.StatusActive},
		{"approved", true},
	}
	sort := mongo.UnorderedDbQuery{
		"id": -1,
	}
	categoryResults := make([]model.BrandInfo, 0)
	if err := p.db.Mongo.List(ctx, p.table, query, skip, limit, &categoryResults, sort); err != nil {
		return nil, err
	}

	return categoryResults, nil
}

// GetBrandDetails ...
func (p *MgoBrand) GetBrandDetails(ctx context.Context, slug string) (*model.BrandInfo, error) {
	q := mongo.DbQuery{
		{"slug", slug},
	}
	brand := &model.BrandInfo{}

	if err := p.db.Mongo.FindOne(ctx, p.table, q, &brand); err != nil {
		return nil, err
	}

	if brand.Status != model.StatusActive {
		return nil, fmt.Errorf("shop is not active")
	}

	if brand.Approved == false {
		return nil, fmt.Errorf("shop is not approved")
	}

	return brand, nil
}

// GetBySlug ...
func (p *MgoBrand) GetBySlug(ctx context.Context, slug string) (*model.BrandInfo, error) {
	q := mongo.DbQuery{
		{"slug", slug},
	}
	brand := &model.BrandInfo{}
	if err := p.db.Mongo.FindOne(ctx, p.table, q, brand); err != nil {
		return nil, err
	}
	return brand, nil
}
