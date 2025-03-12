package author

import (
	"context"
	"errors"
	"fmt"

	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/1755/bookstore-api/internal/pg"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type GetManyParamsSort struct {
	Field     string
	Direction string
}

type GetManyParams struct {
	Offset     uint
	Limit      uint
	Sort       *GetManyParamsSort
	FilterName string
}

type UpdateField func(map[string]any)

type DAO interface {
	GetByID(ctx context.Context, id ID) (*Model, error)
	GetMany(ctx context.Context, params *GetManyParams) ([]*Model, error)
	GetManyByBookID(ctx context.Context, id int32) ([]*Model, error)
	Create(ctx context.Context, model *Model) (*Model, error)
	Update(ctx context.Context, id ID, fields ...UpdateField) (*Model, error)
	Delete(ctx context.Context, id ID) error
}

type BasicDAO struct {
	pool pg.Pool
}

var _ DAO = (*BasicDAO)(nil)

func NewBasicDAO(pool pg.Pool) *BasicDAO {
	return &BasicDAO{pool}
}

func (dao *BasicDAO) GetByID(ctx context.Context, id ID) (*Model, error) {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Select("id", "name", "bio", "created_at", "updated_at").
		From("authors").
		Where(goqu.C("id").Eq(id)).
		Prepared(true)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing get by id query", zap.String("sql", sql), zap.Any("args", args))
	row := dao.pool.QueryRow(ctx, sql, args...)

	var author Model
	if err = row.Scan(&author.ID, &author.Name, &author.Bio, &author.CreatedAt, &author.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundError.New("author not found")
		}
		return nil, err
	}

	return &author, nil
}

func (dao *BasicDAO) GetMany(ctx context.Context, params *GetManyParams) ([]*Model, error) {
	logger := lgr.GetLogger(ctx)

	// pretty ugly query builder, but it works
	query := goqu.Dialect("postgres").
		Select("id", "name", "bio", "created_at", "updated_at").
		From("authors").
		Prepared(true)

	if params.FilterName != "" {
		query = query.Where(goqu.C("name").Like(fmt.Sprintf("%%%s%%", params.FilterName)))
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if params.Sort != nil {
		if params.Sort.Direction == "asc" {
			query = query.Order(goqu.C(params.Sort.Field).Asc())
		} else {
			query = query.Order(goqu.C(params.Sort.Field).Desc())
		}

	} else {
		query = query.Order(goqu.C("id").Asc())
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing get many query", zap.String("sql", sql), zap.Any("args", args))
	rows, err := dao.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, InternalError.Wrap(err, "error executing query")
	}
	defer rows.Close()

	var authors []*Model
	for rows.Next() {
		var author Model
		err := rows.Scan(&author.ID, &author.Name, &author.Bio, &author.CreatedAt, &author.UpdatedAt)
		if err != nil {
			return nil, err
		}
		authors = append(authors, &author)
	}

	return authors, nil
}

func (dao *BasicDAO) Create(ctx context.Context, model *Model) (*Model, error) {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Insert("authors").
		Prepared(true).
		Rows(model).
		Returning("*")

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing create query", zap.String("sql", sql), zap.Any("args", args))
	row := dao.pool.QueryRow(ctx, sql, args...)

	var created Model
	if err := row.Scan(&created.ID, &created.Name, &created.Bio, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}

	return &created, nil
}

func UpdateNameField(name string) UpdateField {
	return func(set map[string]any) {
		set["name"] = name
	}
}

func UpdateBioField(bio string) UpdateField {
	return func(set map[string]any) {
		set["bio"] = bio
	}
}

func (dao *BasicDAO) Update(ctx context.Context, id ID, fields ...UpdateField) (*Model, error) {
	logger := lgr.GetLogger(ctx)

	set := make(map[string]any)
	for _, field := range fields {
		field(set)
	}
	set["updated_at"] = goqu.L("NOW()")

	query := goqu.Dialect("postgres").
		Update("authors").
		Set(goqu.Record(set)).
		Where(goqu.C("id").Eq(id)).
		Returning("*").
		Prepared(true)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing update query", zap.String("sql", sql), zap.Any("args", args))
	row := dao.pool.QueryRow(ctx, sql, args...)

	var updated Model
	err = row.Scan(&updated.ID, &updated.Name, &updated.Bio, &updated.CreatedAt, &updated.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundError.New("author not found")
		}
		return nil, err
	}

	return &updated, nil
}

func (dao *BasicDAO) Delete(ctx context.Context, id ID) error {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Delete("authors").
		Where(goqu.C("id").Eq(id)).
		Prepared(true)

	sql, args, err := query.ToSQL()
	if err != nil {
		return InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing delete query", zap.String("sql", sql), zap.Any("args", args))
	result, err := dao.pool.Exec(ctx, sql, args)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return NotFoundError.New("author not found")
	}

	return nil
}

func (dao *BasicDAO) GetManyByBookID(ctx context.Context, id int32) ([]*Model, error) {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Select("author.id", "author.name", "author.bio", "books.created_at", "books.updated_at").
		From("authors").
		InnerJoin(goqu.T("book_authors"), goqu.On(
			goqu.C("book_authors.author_id").Eq(goqu.T("author.id")),
			goqu.C("book_authors.book_id").Eq(id),
		)).
		Prepared(true)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing get many by author id query", zap.String("sql", sql), zap.Any("args", args))
	rows, err := dao.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Model
	for rows.Next() {
		var book Model
		err := rows.Scan(&book.ID, &book.Name, &book.Bio, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}
