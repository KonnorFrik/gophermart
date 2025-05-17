set -e
addgroup -S appgroup && adduser -S user -G appgroup
go mod tidy
go build -o ../bin/ .

