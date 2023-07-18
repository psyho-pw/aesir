BINARY_NAME=aesir

env-up:
	oci os object put -bn environments --file .env --name ${BINARY_NAME}/.env --no-multipart --force

env-down:
	oci os object get -bn environments --file .env --name ${BINARY_NAME}/.env

run:
	APP_ENV=development go run main.go

compile:
	echo "Compiling for every OS and Platform"
	go build -o out/${BINARY_NAME} main.go
	GOOS=linux GOARCH=arm go build -o out/${BINARY_NAME}-linux-arm main.go
    GOOS=linux GOARCH=arm64 go build -o out/${BINARY_NAME}-linux-arm64 main.go
    GOOS=freebsd GOARCH=386 go build -o out/${BINARY_NAME}-freebsd-386 main.go
	GOOS=windows GOARCH=386 go build -o out/${BINARY_NAME}-windows-386 main.go

build:
	go build -o out/${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=darwin go build -o out/${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o out/${BINARY_NAME}-linux main.go

clean:
	go clean -modcache
	rm -r out

run-prod:
	./out/${BINARY_NAME}

deps:
	go mod download && go mod tidy

vet:
	go vet

docker-build:
	docker build -t ${BINARY_NAME} .

docker-run:
	@if [ !"$$(docker ps -a -q -f name=${BINARY_NAME})" ]; then \
  		if [ "$$(docker ps -aq -f status=exited -f name=${BINARY_NAME})" ]; then \
  			docker rm ${BINARY_NAME}; \
        fi; \
            docker run -it --name ${BINARY_NAME} -p 3000:3000  ${BINARY_NAME}; \
    fi

analysis:
	go build -gcflags '-m=2'

run-air:
	~/golang/bin/air -c .air.toml