SET datestyle = "ISO, DMY";

CREATE EXTENSION IF NOT EXISTS "cube";
CREATE EXTENSION IF NOT EXISTS "earthdistance";
CREATE EXTENSION IF NOT EXISTS "unaccent";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "postgis_topology";

CREATE OR REPLACE FUNCTION to_ascii(bytea, name) RETURNS text STRICT AS 'to_ascii_encname' LANGUAGE internal;
CREATE OR REPLACE FUNCTION lower(text) RETURNS text LANGUAGE internal IMMUTABLE STRICT AS $function$lower$function$;
CREATE OR REPLACE FUNCTION lower(anyrange) RETURNS anyelement LANGUAGE internal IMMUTABLE STRICT AS $function$range_lower$function$;

-- vytvoreni tabulky pro adresni mista primo z CSV
CREATE TABLE import_address_places (
  place_id INTEGER,
  city_id INTEGER,
  city_name VARCHAR(255),
  momc_id INTEGER,
  momc_name VARCHAR(255),
  mop_id INTEGER,
  mop_name VARCHAR(255),
  city_part_id INTEGER,
  city_part_name VARCHAR(255),
  street_id INTEGER,
  street_name VARCHAR(255),
  place_type VARCHAR(30),
  house_number INTEGER,
  street_number INTEGER,
  street_number_char  VARCHAR(10),
  zip  VARCHAR(255),
  x  REAL,
  y REAL,
  valid_from  TIMESTAMP WITHOUT TIME ZONE
);

-- nahrani dat z CSV
COPY import_address_places FROM '/data/UI_ADRESNI_MISTA.csv' (DELIMITER ';', FORMAT CSV, NULL '', ENCODING 'WIN1250');

-- KOD;NAZEV;OBEC_KOD;PLATI_OD;PLATI_DO;DATUM_VZNIKU
CREATE TABLE import_address_cadastral_territories (
  id INTEGER,
  name VARCHAR(255),
  city_id INTEGER,
  valid_from  TIMESTAMP WITHOUT TIME ZONE,
  valid_to  TIMESTAMP WITHOUT TIME ZONE,
  created_at  TIMESTAMP WITHOUT TIME ZONE
);

-- nahrani dat z CSV
COPY import_address_cadastral_territories FROM '/data/UI_KATASTRALNI_UZEMI.csv' (DELIMITER ';', FORMAT CSV, NULL '', ENCODING 'WIN1250');

-- IMPORT ALL TABLES

-- vytvoreni view pro mesta
CREATE MATERIALIZED VIEW view_address_cities AS
  SELECT 
    city_id as id,
    city_name as name,
    regexp_replace(regexp_replace(regexp_replace(lower(to_ascii(convert_to(import_address_places.city_name,'latin2'), 'latin2')), '-', ' '),'\.',' ','g'), '\s+', ' ','g') as name_search
  FROM 
    import_address_places 
  GROUP BY 
    city_id, city_name;

-- vytvoreni view pro mestske casti
CREATE MATERIALIZED VIEW view_address_city_parts AS
  SELECT
    city_part_id as id,
    city_part_name as name,
    regexp_replace(regexp_replace(regexp_replace(lower(to_ascii(convert_to(import_address_places.city_part_name,'latin2'), 'latin2')), '-', ' '),'\.',' ','g'), '\s+', ' ','g') as name_search,
    city_id
  FROM
    import_address_places 
  GROUP BY 
    city_part_id, city_part_name, city_id;

-- vytvoreni view pro ulice
CREATE MATERIALIZED VIEW view_address_streets AS
  SELECT
    street_id as id,
    street_name as name,
    regexp_replace(regexp_replace(regexp_replace(lower(to_ascii(convert_to(import_address_places.street_name,'latin2'), 'latin2')), '-', ' '),'\.',' ','g'), '\s+', ' ','g') as name_search,
    city_id,
    city_part_id,
    zip
  FROM
    import_address_places
  GROUP BY 
    street_id, street_name, city_id, city_part_id, zip;

-- vytvoreni view pro adresni mista
CREATE MATERIALIZED VIEW view_address_places AS
  SELECT
    place_id as id,
    street_id as street_id,
    CASE WHEN place_type='č.p.' THEN NULL
        ELSE house_number
    END as e,
    CASE WHEN place_type='č.p.' THEN house_number
        ELSE NULL
    END as p,
    CASE WHEN street_number_char is NULL THEN street_number::varchar
        ELSE street_number::varchar || street_number_char
    END as o,
    city_id,
    city_part_id,
    zip,
  --  mop_id as mop_name_id,
  --  momc_id as momc_name_id,
    x,
    y,
    ST_X(ST_AsText(ST_Transform(ST_MakePoint(-1 * x,-1 * y), '+init=epsg:5514', '+init=epsg:4326'))) as longitude,
    ST_Y(ST_AsText(ST_Transform(ST_MakePoint(-1 * x,-1 * y), '+init=epsg:5514', '+init=epsg:4326'))) as latitude,
  --  now() as created_at,
  --  now() as updated_at,
    false as imported
  FROM import_address_places;


CREATE MATERIALIZED VIEW view_address_cadastral_territories AS 
SELECT
      import_address_cadastral_territories.id as id,
      import_address_cadastral_territories.name as name,
      view_address_cities.id as city_id,
--      now() as created_at,
--      now() as updated_at,
      regexp_replace(regexp_replace(regexp_replace(lower(to_ascii(convert_to(import_address_cadastral_territories.name,'latin2'), 'latin2')), '-', ' '),'\.',' ','g'), '\s+', ' ','g') as name_search
FROM
      import_address_cadastral_territories
JOIN
      view_address_cities
ON
      view_address_cities.id = import_address_cadastral_territories.city_id;

DROP INDEX IF EXISTS index_address_cadastral_territories_on_id;
CREATE INDEX index_address_cadastral_territories_on_id ON view_address_cadastral_territories USING btree (id);
DROP INDEX IF EXISTS index_address_cities_on_id;
CREATE INDEX index_address_cities_on_id ON view_address_cities USING btree (id);
DROP INDEX IF EXISTS index_address_city_parts_on_id;
CREATE INDEX index_address_city_parts_on_id ON view_address_city_parts USING btree (id);
DROP INDEX IF EXISTS index_address_streets_on_id;
CREATE INDEX index_address_streets_on_id ON view_address_streets USING btree (id);
DROP INDEX IF EXISTS index_address_places_on_id;
CREATE INDEX index_address_places_on_id ON view_address_places USING btree (id);

DROP INDEX IF EXISTS index_address_places_on_street_id;
CREATE INDEX index_address_places_on_street_id ON view_address_places USING btree (street_id);
DROP INDEX IF EXISTS index_address_places_on_city_part_id;
CREATE INDEX index_address_places_on_city_part_id ON view_address_places USING btree (city_part_id);
DROP INDEX IF EXISTS index_address_places_on_city_id;
CREATE INDEX index_address_places_on_city_id ON view_address_places USING btree (city_id);
DROP INDEX IF EXISTS index_address_places_on_lat_lng;
CREATE INDEX  index_address_places_on_lat_lng ON view_address_places USING gist(ll_to_earth(latitude, longitude));
DROP INDEX IF EXISTS index_address_places_on_point;
CREATE INDEX  index_address_places_on_point ON view_address_places USING gist(point(latitude, longitude));


DROP INDEX IF EXISTS index_address_cadastral_territories_on_name_search;
CREATE INDEX index_address_cadastral_territories_on_name_search ON view_address_cadastral_territories USING gin (name_search gin_trgm_ops);
DROP INDEX IF EXISTS index_address_cities_on_name_search;
CREATE INDEX index_address_cities_on_name_search ON view_address_cities USING gin (name_search gin_trgm_ops);
DROP INDEX IF EXISTS index_address_city_parts_on_name_search;
CREATE INDEX index_address_city_parts_on_name_search ON view_address_city_parts USING gin (name_search gin_trgm_ops);
DROP INDEX IF EXISTS index_address_street_names_on_name_search;
CREATE INDEX index_address_street_names_on_name_search ON view_address_streets USING gin (name_search gin_trgm_ops);

SET datestyle = "ISO, MDY";