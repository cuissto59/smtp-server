build:
	@echo "Building ..."
	@go build -o bin src/main.go



run:
	@go run src/main.go



test:
	@echo "Testing ..."
	@go test ./tests/	



clean:
	@echo "Cleaning..."
	@rm -f bin
