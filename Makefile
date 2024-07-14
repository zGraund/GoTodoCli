# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/todoCli/main.go

# Run the application
run:
	@go run cmd/todoCli/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main