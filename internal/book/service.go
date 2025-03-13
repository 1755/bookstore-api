package book

import "context"

type Service interface {
	GetByID(ctx context.Context, id ID) (*Model, error)
	GetMany(ctx context.Context, params *GetManyParams) ([]*Model, error)
	Create(ctx context.Context, model *Model) (*Model, error)
	Update(ctx context.Context, id ID, fields ...UpdateField) (*Model, error)
	Delete(ctx context.Context, id ID) error
	LinkAuthor(ctx context.Context, id ID, authorID int32) error
	UnlinkAuthor(ctx context.Context, id ID, authorID int32) error
}

var _ Service = (*BasicService)(nil)

type BasicService struct {
	dao DAO
}

func NewBasicService(dao DAO) *BasicService {
	return &BasicService{dao}
}

func (s *BasicService) Create(ctx context.Context, model *Model) (*Model, error) {
	return s.dao.Create(ctx, model)
}

func (s *BasicService) Update(ctx context.Context, id ID, fields ...UpdateField) (*Model, error) {
	return s.dao.Update(ctx, id, fields...)
}

func (s *BasicService) Delete(ctx context.Context, id ID) error {
	return s.dao.Delete(ctx, id)
}

func (s *BasicService) GetByID(ctx context.Context, id ID) (*Model, error) {
	return s.dao.GetByID(ctx, id)
}

func (s *BasicService) GetMany(ctx context.Context, params *GetManyParams) ([]*Model, error) {
	return s.dao.GetMany(ctx, params)
}

func (s *BasicService) LinkAuthor(ctx context.Context, id ID, authorID int32) error {
	return s.dao.LinkAuthor(ctx, id, authorID)
}

func (s *BasicService) UnlinkAuthor(ctx context.Context, id ID, authorID int32) error {
	return s.dao.UnlinkAuthor(ctx, id, authorID)
}
