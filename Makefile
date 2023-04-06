BINARY=scanner
OUTPUT=bin

build:
	echo "Building for every OS and Platform"
	go build -o ${OUTPUT}/${BINARY} main.go
	GOOS=freebsd GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-freebsd-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-windows-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o ${OUTPUT}/${BINARY}-darwin-amd64 main.go

run: build
	echo "Preparing for running..."
	mkdir -p tmp/results
	echo "Run the default binary..."
	./${OUTPUT}/${BINARY}

clean:
	go clean
	rm -rf ${OUTPUT}

all: clean build
