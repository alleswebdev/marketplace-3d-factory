FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY .. .
RUN go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" -o 3dfactory ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/3dfactory .

EXPOSE 80
CMD ["./3dfactory"]
