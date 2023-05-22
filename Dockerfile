FROM golang:1.20-alpine as build

ADD . /go/src/mimuw-project
WORKDIR /go/src/mimuw-project

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

COPY --from=build /go/src/mimuw-project/out /app

EXPOSE 8080

ENTRYPOINT ["/app"]
