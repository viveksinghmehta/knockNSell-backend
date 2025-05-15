# Run the application
run:
	@go run cmd/main.go

# Clean the binary and zip file
clean:
	@echo "Cleaning the binary and zip file..."
	@rm -f bootstrap
	@rm -f bootstrap.zip
	@echo "Removed the binary and zip file"

# Build for AWS lambda
build-aws:
	@echo "Creating binary for aws..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap cmd/main.go
	@echo "Giving executable permision to binary"
	chmod +x bootstrap
	@echo "Zipping the binary ..."
	zip bootstrap.zip bootstrap
	@echo "Zip file is created."
