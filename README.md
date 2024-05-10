## Usage

---
### Local
```go run cmd/main.go input.txt```

---
### Containerized
1. ```docker build -t gomputer-club .```
2. ```docker run --rm gomputer-club input.txt```

**NOTA BENE**: Input files should be additionally mounted to the docker container. By default, only input.txt is mounted.

---
