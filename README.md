# Orbi - Technical test

## How to test
###### How to run
Load the environment variables
```shell
export $(grep -v ^# .env.example)
```
Execute the following command
```shell
make run
```
###### How to stop
```shell
make down
```

## Component diagram
The following diagram shows a general idea of the proposed solution.

[See more about user microservice          (A)](a/README.md)

[See more about notification microservice (B)](b/README.md)

![Component diagram](docs/components.png)
## Highlights
- Healthcheck for the user microservice
- The size of each container was optimized
- All microservices implement [Graceful shutdown](https://www.josestg.com/posts/golang/graceful-shutdown-in-go/)
- User microservice implements [The outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html)
- Exposed library to make http calls to the user microservice (reusable code)
- Exposed library to make grpc calls to the notification microservice (reusable code)