# GoFiber Boilerplate

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
    cd src/ && wire
```
### Commands
```bash
    # Clean packages & removed build files
    make clean
    
    # organize dependencies
    make deps
    
    # OS-compatible builds
    make compile
    
    # run application
    make run
    
    # run application via air
    make run-air
```

## Production
```bash
    # Build docker image
    make docker-build
    
    # Run container
    make docker-run
```
