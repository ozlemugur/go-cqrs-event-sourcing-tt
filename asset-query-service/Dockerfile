#  # Asset Query Service (Reads from Query DB and Serves Clients)
# Step 1: Modules caching
FROM golang:1.22-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/asset-query-service ./cmd/app

# Step 3: Final
FROM alpine:latest
ENV TZ=Europe/Istanbul
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone
COPY --from=builder /app/config /config
COPY --from=builder /bin/asset-query-service /asset-query-service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/asset-query-service"]
