1. run in bash: go mod tidy
2. generate a models with 'https://github.com/sqlc-dev/sqlc'
    bash: cd db/sqlc_conf && sqlc generate
2.1 generate a swagger docs
    bash: cd user_service && swag init
3. start docker compose from 'container' folder
4. run migrations with 'https://github.com/golang-migrate/migrate' 
    bash: migrate -path=/db/migrations -database "postgres://admin:admin@0.0.0.0:5432/gophermart?sslmode=disable" up
    fake user for test: login: 'test' password: 'test'
5. run tests from 'api_test' folder or make manual requests

