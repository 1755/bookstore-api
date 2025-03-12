package author

import "time"

type ID int32

type Model struct {
	ID        ID
	Name      string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
