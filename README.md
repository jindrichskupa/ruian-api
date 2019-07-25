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

```json
[
  {
    "id": 9063811,
    "p": 78,
    "zip": "33008",
    "lng": 13.4180709920695,
    "lat": 49.7987100376894,
    "street": "Školní",
    "city": "Zruč-Senec",
    "city_part": "Senec",
    "address_string": "Školní 78, Senec, Zruč-Senec, 33008"
  },
  {
    "id": 9064001,
    "p": 97,
    "zip": "33008",
    "lng": 13.4185320615227,
    "lat": 49.7987472256835,
    "street": "Školní",
    "city": "Zruč-Senec",
    "city_part": "Senec",
    "address_string": "Školní 97, Senec, Zruč-Senec, 33008"
  },
  {
    "id": 9064028,
    "p": 99,
    "zip": "33008",
    "lng": 13.4185282829638,
    "lat": 49.7985229241633,
    "street": "Školní",
    "city": "Zruč-Senec",
    "city_part": "Senec",
    "address_string": "Školní 99, Senec, Zruč-Senec, 33008"
  },
  ...
]
```

* search for address based on GPS coordinates and range

```bash
curl 'localhost:8080/places/search?latitude=49.8009&longitude=13.4193&range=500&limit=3' | jq
```

```json
[
  {
    "id": 9063081,
    "p": 1,
    "zip": "33008",
    "lng": 13.4207954316166,
    "lat": 49.7937274728657,
    "street": "Senecká",
    "city": "Zruč-Senec",
    "city_part": "Senec",
    "address_string": "Senecká 1, Senec, Zruč-Senec, 33008"
  },
  {
    "id": 9063099,
    "p": 2,
    "zip": "33008",
    "lng": 13.4212828255798,
    "lat": 49.7933398312855,
    "street": "Senecká",
    "city": "Zruč-Senec",
    "city_part": "Senec",
    "address_string": "Senecká 2, Senec, Zruč-Senec, 33008"
  },
  {
    "id": 9063102,
    "p": 4,
    "zip": "33008",
    "lng": 13.4201734224913,
    "lat": 49.7928310973284,
    "street": "Strmá",
    "city": "Zruč-Senec",
    "city_part": "Senec",
    "address_string": "Strmá 4, Senec, Zruč-Senec, 33008"
  }
]
```

* list streets based on street and city names name

```bash
curl 'localhost:8080/streets/search?street=Sko&city=Zruč%20-%20S' | jq
```

* list cadastral territories based on name

```bash
curl 'localhost:8080/cadastral_territories/Holesovice' | jq
```

```json
[
  {
    "id": 730122,
    "name": "Holešovice",
    "city": "Praha"
  },
  {
    "id": 641111,
    "name": "Holešovice u Chroustovic",
    "city": "Chroustovice"
  }
]
```
