# # Asset Management Service (Command API - Writes to Kafka)
# Step 1: Modules caching
FROM --platform=linux/amd64 golang:1.22-bullseye as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download && go mod vendor


# Step 2: Builder
FROM --platform=linux/amd64  golang:1.22-bullseye as builder
#RUN apk add --no-cache librdkafka-dev gcc musl-dev pkgconf cyrus-sasl-dev
RUN apt-get update && apt-get install -y gnupg
RUN apt-get update && apt-get install -y gcc libc6-dev librdkafka-dev pkg-config

#RUN apt-get update && apt-get install -y gcc libc6-dev librdkafka-dev pkg-config
COPY --from=modules /go/pkg /go/pkg
COPY . /app
RUN echo "Listing /app directory..." && ls -l /app && ls -l /app/config
WORKDIR /app
RUN CGO_ENABLED=1  GO111MODULE=on GOOS=linux  GOARCH=amd64  \
    go build -mod=vendor -o /bin/asset-management-service  ./cmd/app 
  

# Step 3: Final
FROM --platform=linux/amd64  golang:1.22-bullseye
ENV TZ=Europe/Istanbul
RUN apt-get update && apt-get install -y tzdata
COPY --from=builder /app/config /config
COPY --from=builder /bin/asset-management-service /asset-management-service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN pwd
WORKDIR /
RUN echo "Listing /app directory..." && ls -l /config && ls -l 

RUN ls -l /config && ls -l /asset-management-service

CMD ["/asset-management-service"]




