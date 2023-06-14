run d:
	go run cmd/shortener/main.go -d "postgres://pavelryltsov:1638@localhost:5432/testdb?sslmode="
run f:
	go run cmd/shortener/main.go -f "memory.log"
run:
	go run cmd/shortener/main.go


