package book

import (
	"context"
	"errors"
	"fmt"

	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/1755/bookstore-api/internal/pg"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type GetManyParamsSort struct {
	Field     string
	Direction string
}

type GetManyParams struct {
	Offset      uint
	Limit       uint
	Sort        *GetManyParamsSort
	FilterTitle string
}

type UpdateField func(map[string]any)

type DAO interface {
	GetByID(ctx context.Context, id ID) (*Model, error)
	GetMany(ctx context.Context, params *GetManyParams) ([]*Model, error)
	GetManyByAuthorID(ctx context.Context, id int32) ([]*Model, error)
	Create(ctx context.Context, model *Model) (*Model, error)
	Update(ctx context.Context, id ID, fields ...UpdateField) (*Model, error)
	Delete(ctx context.Context, id ID) error
	LinkAuthor(ctx context.Context, id ID, authorID int32) error
	UnlinkAuthor(ctx context.Context, id ID, authorID int32) error
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
		Select("id", "title", "summary", "published_year", "created_at", "updated_at").
		From("books").
		Where(goqu.C("id").Eq(id)).
		Prepared(true)

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing get by id query", zap.String("sql", sql), zap.Any("args", args))
	row := dao.pool.QueryRow(ctx, sql, args...)

	var book Model
	if err := row.Scan(&book.ID, &book.Title, &book.Summary, &book.PublishedYear, &book.CreatedAt, &book.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundError.New("book not found")
		}
		return nil, err
	}

	return &book, nil
}

func (dao *BasicDAO) GetMany(ctx context.Context, params *GetManyParams) ([]*Model, error) {
	logger := lgr.GetLogger(ctx)

	// pretty ugly query builder, but it works
	query := goqu.Dialect("postgres").
		Select("id", "title", "summary", "published_year", "created_at", "updated_at").
		From("books").
		Prepared(true)

	if params.FilterTitle != "" {
		query = query.Where(goqu.C("name").Like(fmt.Sprintf("%%%s%%", params.FilterTitle)))
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
		return nil, err
	}
	defer rows.Close()

	var books []*Model
	for rows.Next() {
		var book Model
		err := rows.Scan(&book.ID, &book.Title, &book.Summary, &book.PublishedYear, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (dao *BasicDAO) Create(ctx context.Context, model *Model) (*Model, error) {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Insert("books").
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
	if err := row.Scan(&created.ID, &created.Title, &created.Summary, &created.PublishedYear, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}

	return &created, nil
}

func UpdateTitleField(title string) UpdateField {
	return func(set map[string]any) {
		set["title"] = title
	}
}

func UpdateSummaryField(summary string) UpdateField {
	return func(set map[string]any) {
		set["summary"] = summary
	}
}

func UpdatePublishedYearField(publishedYear int32) UpdateField {
	return func(set map[string]any) {
		set["published_year"] = publishedYear
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
		Update("books").
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
	if err := row.Scan(&updated.ID, &updated.Title, &updated.Summary, &updated.PublishedYear, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundError.New("books not found")
		}
		return nil, err
	}

	return &updated, nil
}

func (dao *BasicDAO) Delete(ctx context.Context, id ID) error {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Delete("books").
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
		return NotFoundError.New("book not found")
	}

	return nil
}

func (dao *BasicDAO) GetManyByAuthorID(ctx context.Context, id int32) ([]*Model, error) {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Select("books.id", "books.title", "books.summary", "books.published_year", "books.created_at", "books.updated_at").
		From("books").
		InnerJoin(goqu.T("book_authors"), goqu.On(
			goqu.C("book_authors.book_id").Eq(goqu.T("books.id")),
			goqu.C("book_authors.author_id").Eq(id),
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
		err := rows.Scan(&book.ID, &book.Title, &book.Summary, &book.PublishedYear, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (dao *BasicDAO) LinkAuthor(ctx context.Context, id ID, authorID int32) error {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Insert("book_authors").
		Prepared(true).
		Rows(map[string]any{"book_id": id, "author_id": authorID})

	sql, args, err := query.ToSQL()
	if err != nil {
		return InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing link author query", zap.String("sql", sql), zap.Any("args", args))
	_, err = dao.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (dao *BasicDAO) UnlinkAuthor(ctx context.Context, id ID, authorID int32) error {
	logger := lgr.GetLogger(ctx)

	query := goqu.Dialect("postgres").
		Delete("book_authors").
		Where(goqu.C("book_id").Eq(id)).
		Where(goqu.C("author_id").Eq(authorID))

	sql, args, err := query.ToSQL()
	if err != nil {
		return InternalError.Wrap(err, "error building query")
	}

	logger.Debug("executing unlink author query", zap.String("sql", sql), zap.Any("args", args))
	result, err := dao.pool.Exec(ctx, sql)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return NotFoundError.New("link not found")
	}

	return nil
}
