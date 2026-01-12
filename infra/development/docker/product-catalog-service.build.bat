set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o build/product-catalog-service ./services/product-catalog-service/cmd/main.go