# decompose
Convert Docker Compose YAML into Docker commands.

This package is still in flux. And has bugs.

## Installation
Install or upgrade Decompose with this command.
```
go get -u github.com/variadico/decompose/cmd/decompose
```

## Usage
```
decompose docker-compose.yml
```

If your `docker-compose.yml` looks like this,
```
postgres:
    image: postgres:latest
    ports:
        - "5432"

myapp:
    build: .
    links:
        - postgres:db
    ports:
        - "80:80"
    environment:
        - FOO=bar
        - FIZZ=buzz
```

then, `decompose` will output these commands.
```
docker run --name=dir_postgres --publish=5432 postgres:latest
docker build --tag dir_myapp .
docker run --name=dir_myapp --env=FOO=bar --env=FIZZ=buzz --link=postgres:db --publish=80:80 dir_myapp
```

Names will be prefixed with the name of the directory that `docker-compose.yml`
is in.
