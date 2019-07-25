# RUIAN API Server

## Description

Public REST API for Czech RUIAN Database based on CSV exports

### TODO

* [ ] Automate load from CSV
* [ ] Automate update from CSV
* [ ] Generate exports - full (CSV, JSON, SQLite)
* [ ] Generate exports - diff (CSV, JSON, SQLite)
* [ ] Improvements...

## Build

### Download dependencies

```bash
dep ensure
```

### Make

```bash
make

# or cross-compile for Linux amd64
make linux64

# or docker image
make docker
```

## Run

### Configuration ENV variables

Variables and default values

```bash
RUIAN_LISTEN_IP=0.0.0.0
RUIAN_LISTEN_PORT=8080
RUIAN_DB_HOSTNAME=localhost
RUIAN_DB_PORT=5432
RUIAN_DB_USER=postgres
RUIAN_DB_PASSWORD=password
RUIAN_DB_NAME=ruiandb
RUIAN_DB_RETRIES=10
```

### Docker compose

1. update `docker-compose.yml` environment variables values
2. run `docker-compose up -d`

### Standalone

```bash
RUIAN_DB_NAME=ruian_db_test ./ruian-api
```

## Database setup

### Requirements

1. PostgreSQL + PostGIS
2. Basic PostgreSQL extenstions

### Setup

1. download ziped CSV files
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

### API Call Examples

* search for address based on part of address

```bash
curl 'localhost:8080/places/search?street=Sko&city=Zruč%20-%20S&city_part=Senec' | jq
```

* search for address based on GPS coordinates and range

```bash
curl 'localhost:8080/places/search?latitude=49.8009&longitude=13.4193&range=500&limit=10' | jq
```

* list streets based on street and city names name

```bash
curl 'localhost:8080/streets/search?street=Sko&city=Zruč%20-%20S' | jq
```
