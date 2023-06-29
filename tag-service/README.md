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

### Options

There are several options for the service, which can be set as environment variables in a `.env` file or directly in a docker container. When left empty, they will reset to a default value.

- MONGO_HOST: IPV4 of MongoDB access point (default: localhost)
- MONGO_PORT: Port where MongoDB listens (default: 27017)
- PORT: Port on which the serivce runs (default: 3000)
- DEBUG: Prints debug information in the log and saves generated and true results of queries in the database (default: 0)
- AGGR: This flag activates whether use case 3 aggregations are calculated or deactivated and returning a 501 not implemented code. (default: 1)
