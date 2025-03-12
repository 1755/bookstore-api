package book

import "time"

type ID int32

type Model struct {
	ID            ID
	Title         string
	Summary       string
	PublishedYear int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
