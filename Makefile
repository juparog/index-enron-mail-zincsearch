clean:
	echo "cleaning..."
	rm -r -f ./build

install:
	echo "installing..."
	go mod tidy

compile:
	echo "compiling..."
	GOOS=linux go build -o ./build/main-linux-386/main ./cmd/lambda/main.go

package:
	echo "packing..."
	cd ./build/main-linux-386/ && \
	zip main.zip main

all: clean install compile package
