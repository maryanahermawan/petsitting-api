## Services
Will run Postgres on default port with default config for local
```shell
docker-compose up --build
```

## Database migrations
```shell
DATABASE_PASSWORD="password" go run cmd/migrate/main.go
```

To run a down migration
```shell
DATABASE_PASSWORD="password" go run cmd/migrate/main.go -direction=down
```

## Build and run app via Docker
```shell
docker build . -f ./docker/Dockerfile -t golang-starter:latest
docker run -p 8000:8000 golang-starter:latest
```
