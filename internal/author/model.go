package author

import "time"

type ID int32

type Model struct {
	ID        ID        `db:"id"`
	Name      string    `db:"name"`
	Bio       string    `db:"bio"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
