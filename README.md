# Aesir - inactivity monitor slackbot

## Development
### Start application
- Live reload by [Air](https://github.com/cosmtrek/air)
```bash
    make deps
    make run or make run-air
```
### Require [Google Wire](https://github.com/google/wire) for DI
```bash
    # go install github.com/google/wire/cmd/wire@latest
    # run following command from project root
    wire . or make generate
```
### Commands
```bash
    # Clean packages & removed build files
    make clean
    
    # organize dependencies
    make deps

    # tidy dependencies
    make tidy // careful!; removes wire subcommands

    # clean mod cache and build files
    make clean
    
    # OS-compatible builds
    make compile
    
    # run application
    make run
    
    # run application via air
    make run-air
```

## Test
### Test and mocking by [Testify](https://github.com/stretchr/testify) & [Mockery](https://github.com/vektra/mockery)
```bash
    # brew install mockery
    # run following command from project root
    make generate
    make test
```

## Production
```bash
    # Build docker image
    make docker-build
    
    # Run container
    make docker-run
```
