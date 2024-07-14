# database-sharding-poc
## Introduction

This project is a simple proof of concept to demonstrate the concepts of "database sharding" and "consistent hashing". It provides a basic implementation to help understand how data can be distributed across multiple database instances using consistent hashing to achieve sharding.

## Docker

- Build the docker image:  
```bash
docker build -t postgres-shard .
```

- Run the containers:  
```bash
docker run --name shard-1 -p 5431:5432 -d -e POSTGRES_DB=test_db -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test postgres-shard
docker run --name shard-2 -p 5433:5432 -d -e POSTGRES_DB=test_db -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test postgres-shard
docker run --name shard-3 -p 5434:5432 -d -e POSTGRES_DB=test_db -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test postgres-shard
```

- Info: You can use [pgAdmin](https://www.pgadmin.org/download/) to make the connections to the different shards.

## Golang
- Setup Golang dependencies:  
```bash
cd app
go mod tidy
```

- Run application:  
```bash
cd app
go run .
```

## CURLs

#### POST
```bash
curl -X POST "http://localhost:8081/api/products/post?product_name=laptop"

```
On windows
```bash
Invoke-WebRequest -Uri "http://localhost:8081/api/products/post?product_name=laptop" -Method POST
```

#### GET
```bash
curl "http://localhost:8081/api/products/get?product_code=abcde"
```
On windows
```bash
Invoke-WebRequest -Uri "http://localhost:8081/api/products/get?product_code=abcde" -Method GET
```
