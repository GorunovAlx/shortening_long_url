module github.com/GorunovAlx/shortening_long_url

go 1.17

replace github.com/GorunovAlx/shortening_long_url/internal/app/handlers => ../internal/app/handlers

replace github.com/GorunovAlx/shortening_long_url/internal/app/storage => ../internal/app/storage

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
