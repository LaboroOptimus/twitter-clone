module user

go 1.24.4

require (
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgx/v5 v5.7.5
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.41.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

require (
	github.com/redis/go-redis/v9 v9.13.0
	twitter/pkg/redisx v0.0.0-00010101000000-000000000000
)

replace twitter/pkg/redisx => ../../pkg/redisx
