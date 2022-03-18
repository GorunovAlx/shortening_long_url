module github.com/GorunovAlx/shortening_long_url

go 1.17

replace github.com/GorunovAlx/shortening_long_url/internal/app/handlers => ../internal/app/handlers

replace github.com/GorunovAlx/shortening_long_url/internal/app/storage => ../internal/app/storage

replace github.com/GorunovAlx/shortening_long_url/internal/app/configs => ../internal/app/configs

replace github.com/GorunovAlx/shortening_long_url/internal/app/generators => ../internal/app/generators

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/caarlos0/env/v6 v6.9.1
	github.com/go-chi/chi/v5 v5.0.7
	github.com/itchyny/base58-go v0.2.0
	github.com/jackc/pgx/v4 v4.15.0
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.11.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.10.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
