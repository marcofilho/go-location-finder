FROM golang:latest AS builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-w -s" -o go-location-finder main.go

FROM alpine:latest
COPY --from=builder /app /app
CMD ["/app/go-location-finder"]