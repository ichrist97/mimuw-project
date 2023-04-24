FROM golang:1.20-alpine as build

WORKDIR /app

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

COPY --from=build /app/out /app

EXPOSE 8080

ENTRYPOINT ["/app"]
