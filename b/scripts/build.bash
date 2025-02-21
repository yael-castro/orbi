#!/bin/bash

# Variables for only read
runtime="github.com/yael-castro/orbi/b/internal/runtime"
# commit=$(git log --pretty=format:'%h' -n 1)
commit='unknown'

# Command arguments
subcommand="$1"
shift

ldflags=""
options=""

function build() {
    cd "./cmd/$binary" || exit

    if ! go mod tidy
    then
      exit 1
    fi

    if ! go build \
      -o ../../build/ \
      -tags "$tags" \
      -ldflags "$ldflags" \
      "$options"
    then
      exit 1
    fi

    cd ../../

    echo "MD5 checksum: $(md5sum "build/$binary")"
    echo "Success build"
    exit
}


if [ "$subcommand" = "consumer" ]; then
  ldflags='-extldflags "-static" -linkmode external -w -s '
  binary="notifications-consumer"
  tags="consumer,musl"

  printf "\nBuilding CLI in \"build\" directory\n"
  CGO_ENABLED=1 build
fi

if [ "$subcommand" = "grpc" ]; then
  binary="notifications-grpc"
  tags="grpc"

  printf "\nBuilding gRPC server in \"build\" directory\n"
  CGO_ENABLED=0 build
fi

echo "Invalid subcommand: $subcommand"
exit 1