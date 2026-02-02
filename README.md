# Cineguard

A nice tool to manage your favorite movies or those you wanna watch

## Quickstart

Build :

```sh
make build 
```

Start a PostgreSql database :

```sh
podman run --name cineguard-db \
        -e POSTGRES_PASSWORD=mypassword \
        -e POSTGRES_USER=myuser \
        -e POSTGRES_DB=cineguard \
        -p 5432:5432 \
        -v cineguard-data:/var/lib/postgresql \
        docker.io/library/postgres
```


```sh
./cineguard serve
# ...
# [GIN-debug] Listening and serving HTTP on 127.0.0.1:8080
```

Check status :

```sh
curl localhost:8080/api/v1/health
```

## Build

```sh
make build
```

## TODO

models for sql tables have been done.
Now it is time to create the database, the tables and start to experiment with it.
The post method for movie is implemented, now test it.