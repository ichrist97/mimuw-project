# Tag service

Web service for handling the creation of user_tags, querying of user_profiles and querying aggregations.

Built using the [Go Fiber](https://gofiber.io/) framework.

## Setup

### Local

Download dependencies:

```
go mod download
```

Run app:

```
go run server.go
```

### Docker

Build image:

```
docker build -t tag-service .
```

Run container:

```
docker run -p 3000:3000 -e MONGO_HOST=<DB_IP> -e DEBUG=<DEBUG_MODE> tag-service
```
