FROM golang:1.20-alpine as build

ADD . /go/src/aggregation-service
WORKDIR /go/src/aggregation-service

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

# build app
RUN go build -o ./out .

# Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /go/src/aggregation-service/out /app

EXPOSE 8080

ENTRYPOINT ["/app"]
