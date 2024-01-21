#! /bin/bash

if ! [ -x "$(command -v go)" ]; then
  echo 'Error: go is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v swag)" ]; then
  echo 'Error: swag docs is not installed.' >&2
  echo 'Please install it with the following command:'
  echo 'go install github.com/swaggo/swag/cmd/swag@latest'
  exit 1
fi

if ! [ -x "$(command -v docker)" ]; then
  echo 'Error: docker is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v docker-compose)" ]; then
  echo 'Error: docker-compose is not installed.' >&2
  exit 1
fi

if  [ -z "$1" ] ; then
    echo "Error: Please specify which binary under cmd/ you want to run."
    echo "Example: ./start.sh main.go dev"
    exit 1
fi

if [[ "$1" == "swagger" ]]; then
    echo "Running re-building swagger docs" 
      swag init --generalInfo core/services/api/docs/docs.go  --parseDependency --generatedTime -o core/services/api/docs;
     exit 0
fi

echo 'Running dev environment, if you want to run the production one do'
echo 'start.sh binary prod'
echo 'If you want to run tests and quit, do'
echo 'start.sh binary test'
echo ''
echo '---------------------------------------------------------------------------'

if [[ "$1" == "test" ]] ; then
    echo "Running dev environment"
    docker-compose -f scripts/dev.docker-compose.yml up -d --build
    go mod download
    echo "Running tests ...." 
    go test ./core
    go test ./cmd
    go test ./pkgs
    go test ./utils
    docker-compose down
    exit 0
elif [[ "$1" == "prod" ]]; then
    echo "Running prod environment"
    docker-compose -f ./scripts/docker-compose.yml  up -d --build
elif [[ "$1" == "kill" ]]; then 
    docker-compose -f ./scripts/docker-compose.yml down 
    docker-compose -f ./scripts/dev.docker-compose.yml down 
else
    echo "Running dev environment"
    docker-compose -f scripts/dev.docker-compose.yml up -d --build
    go mod download
    go run "cmd/inariam/main.go" $2
fi