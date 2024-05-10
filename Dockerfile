FROM golang:1.22 as BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /gomputerClub cmd/main.go

EXPOSE 8080

FROM alpine:latest

WORKDIR /
COPY --from=BuildStage /gomputerClub /gomputerClub
COPY entrypoint.sh /app/entrypoint.sh
COPY input.txt /input.txt

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
