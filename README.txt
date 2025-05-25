1. start docker compose from 'container' folder
2. generate a models with 'https://github.com/sqlc-dev/sqlc'
    bash: cd db/sqlc_conf && sqlc generate
3. run migrations with 'https://github.com/golang-migrate/migrate' 
    bash: migrate -path=/db/migrations -database "postgres://admin:admin@0.0.0.0:5432/gophermart?sslmode=disable" up
4. run tests from 'api_test' folder or make manual requests
