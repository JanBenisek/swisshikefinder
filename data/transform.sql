/* ============== TRANSFORM ============== 
 * transforms nested data into semi-normalised tables
 * some tables still need further normalisation with:
 *  unnest(additionalDescriptions, recursive := true) 
 * TO-DO: special schema and user?
 * */

create schema if not exists bronze;
create schema if not exists silver;
create schema if not exists gold;
create schema if not exists recommendations;

create view if not exists bronze.tours as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/raw_data/tours_data.json');
-- create view if not exists bronze.destinations as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/raw_data/destinations_data.json');
-- create view if not exists bronze.attractions as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/raw_data/attractions_data.json');



/*
drop view bronze.tours;
drop view bronze.destinations;
drop view bronze.attractions;

drop schema if exists bronze cascade;
drop schema if exists silver cascade;
drop schema if exists gold cascade;

-- TOURS
drop table gold.tour_tourist_types;
drop table gold.tour_itinerary;
drop table gold.tour_images;
drop table gold.tour_classification;
drop table gold.tours;

select * from bronze.tours limit 10;

select * from gold.tour_tourist_types limit 10;
select * from gold.tour_itinerary limit 10;
select * from gold.tour_images limit 10;
select * from gold.tour_classification limit 10;
select * from gold.tours limit 10;

*/


/* ============== TOURS ============== */
-- Tourist types of Tour, eg. Snow Lover, Outdoor Enthusiast - Biker and Cyclist, ...
-- ID 1:m tourist_type
create or replace table gold.tour_tourist_types as
select
    identifier as ID,
    unnest(touristType) as tourist_type
from bronze.tours;


-- Itinerary of Tour, eg. Tanzboden -> Chrüzegg, ...
-- ID 1:m name
create or replace table gold.tour_itinerary as
with btbl as (
    select
        identifier as ID,
        unnest(itinerary, recursive := true) as itinerary
    from bronze.tours
    )
select
    ID,
    "@type" as type,
    identifier as itinerary_id,
    name,
    url
from btbl;


-- Images of Tour, eg. url1, url2, ...
-- ID 1:m url
create or replace table gold.tour_images as
with btbl as (
    select
        identifier as ID,
        unnest(image, recursive := true) as image
    from bronze.tours
    )
select
    ID,
    "@type" as type,
    keywords,
    publisher,
    encodingFormat as encoding_format,
    url,
    name,
    copyrightHolder as copyright_holder,
    width,
    height
from btbl;


-- Classification of Tour, eg. National, At the lake, ...
-- Note: might be usefu to filter by `name` (landscape, views, reachabilitylocation, ...)
-- ID 1:m name/classification
create or replace table gold.tour_classification as
with btbl as (
    select
        identifier as ID,
        unnest(classification, recursive := true) as image
    from bronze.tours
    ),
explode_values as (
    select
        ID,
        "@context" as context,
        name,
        title,
        unnest(values, recursive := true) as vals
    from btbl
    )
select
    ID,
    context,
    name,
    title_1 as classification
    
from explode_values
where name_1 != 'bysa';


-- Tour, eg. National, At the lake, ...
-- ID 1:1 "attributes in the table"
create or replace table gold.tours as
with btbl as (
    select
        identifier as ID,
        "@type" as record_type,
        name,
        abstract,
        partOfTrip.identifier as part_of_trip_id,
        url as url_mysw,
        geo.latitude as lat,
        geo.longitude as lon,
        geo.elevation as elevation,
        -- switzerland mobility info
        switzerlandMobility.logo as logo_url,
        switzerlandMobility.stage as stage,
        switzerlandMobility.requirements.technical as requirements_technical,
        switzerlandMobility.requirements.endurance as requirements_endurance,
        switzerlandMobility.routeCategory as route_category,
        switzerlandMobility.url as url_swmo,
        -- specifications
        specs.season as season,
        specs.distance as distance,
        specs.duration as duration,
        specs.durationReverseDirection as duration_reverse_direction,
        specs.ascent as ascent,
        specs.descent as descent,
        specs.barrierFree as barrier_free
        
    from bronze.tours
    )
select
    b.*
from btbl b;


/* ============== RECOMMENDATIONS ============== */
--drop table recommendations.tours
CREATE SEQUENCE rec_id_seq START 1;
CREATE TABLE recommendations.tours(id INTEGER DEFAULT nextval('rec_id_seq'), "title" VARCHAR, "description" VARCHAR, created_at TIMESTAMP, expires_at TIMESTAMP);

--insert into recommendations.tours (title, description, created_at, expires_at) 
--values('foo', 'long description', current_timestamp, current_timestamp + interval 3 day) RETURNING id;

--select * from recommendations.tours


/* ============== EXPPORT DB ============== */
--EXPORT DATABASE '/Users/janbenisek/GithubRepos/swisshikefinder/data/db' (
EXPORT DATABASE '/Users/janbenisek/GithubRepos/swisshikefinder/data/db2' (
    FORMAT PARQUET,
    COMPRESSION ZSTD
);


/* ============== EXPPORT DB ============== */
IMPORT DATABASE '/Users/janbenisek/GithubRepos/swisshikefinder/data/db2';
