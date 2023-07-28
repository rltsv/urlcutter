run d:
	go run cmd/shortener/main.go -d "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"
run f:
	go run cmd/shortener/main.go -f "memory.log"
run:
	go run cmd/shortener/main.go


