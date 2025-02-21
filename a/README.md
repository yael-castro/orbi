# Users microservice

This project provides a REST API for manage operations related to users (create, update and get)
in which I implement the [transactional outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html)
with [the hexagonal pattern](https://alistair.cockburn.us/hexagonal-architecture/) to have the `service/sender` and
the `message relay` in the same code base to avoid having different projects sharing databases.

![Component diagram](docs/images/components.svg)

As the diagram shows, this project is compiled in two parts:
1. `users-http`  It is the binary in charge of manage the http requests related with users (service/sender)
2. `users-relay` It is the binary in charge of read and send the messages (message relay)

### Getting started
[See the required environment variables](.env.example)

[See the OpenAPI specification](docs/OpenAPI.json)

###### Database schema
The database schema is described by the `.sql` files in the [sql](scripts/sql) directory.

### How to use from source
Follow the instructions below to compile and locally.
(All compiled binaries will put in the `build` directory)

> These binaries require external dependencies to run (a PostgreSQL server with a specific database schema)
> so should be run using docker compose. [See documentation](../README.md)
###### User service (sender)
```shell
make http
./build/users-http
```
###### Message relay
```shell
make relay
./build/users-relay
```

### Architecture decisions
###### Go project layout standard
I decided to follow the [Go project layout standard](https://github.com/golang-standards/project-layout).
###### Package tree
I built the package tree following the concepts of the [hexagonal architecture pattern](https://alistair.cockburn.us/hexagonal-architecture/).
```
.
├── cmd
├── internal
│   ├── app
│   │   ├── business (Use cases, business rules, data models and ports)
│   │   ├── input    (Everything related to "drive" adapters)
│   │   └── output   (Everything related to "driven" adapters)
│   └── container (DI container)
└── pkg (Public and global code, potencially libraries)
```
###### Compile only what is required
According to the theory of hexagonal architecture, it is possible to have *n* adapters for different external signals (http, gRPC, command line).

So I decided to compile a binary to handle each signal.