# RUIAN API Server

## Description

## Build

```bash
make

# or cross-compile for Linux amd64
make linux64

# or docker image
make docker
```

## Run

### Docker compose

1. update `docker-compose.yml` environment variables values
2. run `docker-compose up -d`

### Standalone

```bash
RUIAN_DB_HOSTNAME=localhost RUIAN_DB_PORT=5432 RUIAN_DB_USER=postgres RUIAN_DB_PASSWORD=password RUIAN_DB_NAME=ruian ./ruian-api
```

## Database setup

### Requirements

1. PostgreSQL + PostGIS
2. Basic PostgreSQL extenstions

### Setup

1. download data ziped CSV files
2. extract .zip files
3. join files and delete CSV headers with scripts
4. update paths in `import.sql` script
5. run `import.sql` script

### Data URLs

* Address Places:
* Cadastral Teritories:

## Usage

### Rest API

* [openapi.yml](./openapi.yml)

### API Call examples

* search for address based on part of address

```bash
curl 'localhost:8080/places/search?street=Sko&city=Zruč%20-%20S&city_part=Senec' | jq
```

* list streets based on name

```bash
curl 'localhost:8080/streets/Skolni' | jq
```
