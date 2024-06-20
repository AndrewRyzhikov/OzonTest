FROM mirror.gcr.io/golang:1.22-alpine as builder

WORKDIR /src
COPY go.* ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/
FROM mirror.gcr.io/alpine

WORKDIR /app
COPY --from=builder /src/app .
CMD ["./app"]
