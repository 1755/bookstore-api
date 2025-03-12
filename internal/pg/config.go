package pg

type Config struct {
	URI      string `validate:"required,url"`
	MinConns int32  `validate:"min=0"`
	MaxConns int32  `validate:"required,min=1"`
}
