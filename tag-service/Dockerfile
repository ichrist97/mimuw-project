FROM golang:1.20-alpine as build

ADD . /go/src/tag-service
WORKDIR /go/src/tag-service

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

COPY --from=build /go/src/tag-service/out /app

EXPOSE 8080

ENTRYPOINT ["/app"]
