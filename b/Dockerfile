FROM golang:tip-alpine3.21 AS builder

WORKDIR /app

# Install "make" command
RUN apk add --no-cache make

# Install GCC
RUN apk add --no-cache build-base

# Git installation:
#   Git is required to extract the hash of the last commit to use it as the version of the compiled binaries.
#   TODO: Because the .git directory is not here . because this repository contains two projects it is required to set additional_context.
#   TODO: In a normal scenario where there is only one project per repository this would not be a problem.
RUN apk add --no-cache git

COPY . .

# Compiling only whats is required
ARG TAG
RUN make "$TAG"

FROM alpine:3.21.3

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/build .

# Executing compiled binary
CMD "./$EXECUTABLE"