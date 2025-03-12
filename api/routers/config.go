package routers

type Config struct {
	BaseURL string `validate:"required,url"`
}
