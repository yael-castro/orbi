# Orbi - Technical test


## How to run
Follow the instructions below to up the containers.
```shell
export $(grep -v ^# .env.example)
make run
```
If you want stop and clean all execute the command below.
```shell
make down
```

## Component diagram
The following diagram shows a general idea of the proposed solution.

[See more about user microservice          (A)](a/README.md)

[See more about notifications microservice (B)](b/README.md)

![Component diagram](docs/components.png)

## Highlights
- All services implement [Graceful shutdown]()
- User microservices implements [The outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html)
- The size of each container was optimized