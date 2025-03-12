package book

import "time"

type ID int32

type Model struct {
	ID            ID        `db:"id"`
	Title         string    `db:"title"`
	Summary       string    `db:"summary"`
	PublishedYear int32     `db:"published_year"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
