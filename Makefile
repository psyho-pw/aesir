BINARY_NAME=aesir

env-up:
	oci os object put -bn environments --file .env/.env.development --name ${BINARY_NAME}/.env/.env.development --no-multipart --force
	oci os object put -bn environments --file .env/.env.production --name ${BINARY_NAME}/.env/.env.production --no-multipart --force
	oci os object put -bn environments --file .env/aesir.json --name ${BINARY_NAME}/.env/aesir.json --no-multipart --force

env-down:
	oci os object get -bn environments --file .env/.env.development --name ${BINARY_NAME}/.env/.env.development
	oci os object get -bn environments --file .env/.env.production --name ${BINARY_NAME}/.env/.env.production
	oci os object get -bn environments --file .env/aesir.json --name ${BINARY_NAME}/.env/aesir.json

run:
	APP_ENV=development go run main.go

compile:
	echo "Compiling for every OS and Platform"
	go build -o out/${BINARY_NAME} main.go wire_gen.go
	GOOS=linux GOARCH=arm go build -o out/${BINARY_NAME}-linux-arm main.go wire_gen.go
    GOOS=linux GOARCH=arm64 go build -o out/${BINARY_NAME}-linux-arm64 main.go wire_gen.go
    GOOS=freebsd GOARCH=386 go build -o out/${BINARY_NAME}-freebsd-386 main.go wire_gen.go
	GOOS=windows GOARCH=386 go build -o out/${BINARY_NAME}-windows-386 main.go wire_gen.go
	GOARCH=amd64 GOOS=darwin go build -o out/${BINARY_NAME}-darwin main.go wire_gen.go
    GOARCH=amd64 GOOS=linux go build -o out/${BINARY_NAME}-linux main.go wire_gen.go

build:
	go build -o out/${BINARY_NAME} .

deps:
	go mod download

tidy:
	go mod tidy

clean:
	go clean -modcache
	rm -r out

generate:
	go generate ./...

test:
	APP_ENV=development go test -v ./...

run-prod:
	./out/${BINARY_NAME}

vet:
	go vet

docker-build:
	sudo docker build --platform linux/amd64 --tag fishcreek/${BINARY_NAME} -f Dockerfile .

docker-push:
	sudo docker push fishcreek/${BINARY_NAME}

docker-run:
	@if [ !"$$(docker ps -a -q -f name=${BINARY_NAME})" ]; then \
  		if [ "$$(docker ps -aq -f status=exited -f name=${BINARY_NAME})" ]; then \
  			docker rm ${BINARY_NAME}; \
        fi; \
            docker run -it --name ${BINARY_NAME} -p 8000:8000  fishcreek/${BINARY_NAME}; \
    fi

analysis:
	go build -gcflags '-m=2'

run-air:
	air -c .air.toml