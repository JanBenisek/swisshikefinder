/* ============== TRANSFORM ============== 
 * transforms nested data into semi-normalised tables
 * some tables still need further normalisation with:
 *  unnest(additionalDescriptions, recursive := true) 
 * TO-DO: special schema and user?
 * */

create schema if not exists bronze;
create schema if not exists silver;
create schema if not exists gold;

create view if not exists bronze.tours as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/tours_data.json');
create view if not exists bronze.destinations as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/destinations_data.json');
create view if not exists bronze.attractions as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/attractions_data.json');


/*
drop view tours_raw;
drop view destinations_raw;
drop view attractions_raw;

drop table tours;
drop table destinations;
drop table attractions;

select * from tours limit 10;
select * from destinations limit 10;
select * from attractions limit 10;

select * from tours_raw limit 10;
select * from destinations_raw limit 10;
select * from attractions_raw limit 10;

*/


/* ============== TOURS ============== */
-- Tourist types of Tour, eg. Snow Lover, Outdoor Enthusiast - Biker and Cyclist, ...
-- ID 1:m tourist_type
create or replace table gold.tour_tourist_types
select
    identifier as ID,
    unnest(touristType) as tourist_type
from bronze.tours;


-- Itinerary of Tour, eg. Tanzboden -> Chrüzegg, ...
-- ID 1:m name
create or replace table gold.tour_itinerary
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
create or replace table gold.tour_images
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
create or replace table gold.tour_classification
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
with 
btbl as (
    select
        identifier as id,
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
from btbl b










describe gold.tours;











/* ============== DESTINATIONS ============== */
create or replace table  destinations as
with 
btbl as (
    select
        identifier as id,
        "@type" as record_type,
        name,
        category,
        abstract,
        description,
        photo,
        url as url_mysw,
        additionalDescriptions as additional_description, -- into a separate table
        address, -- also into separate table later
        geo.latitude as lat,
        geo.longitude as lon,
        image, -- list of images
        containedInPlace as contained_in_destination -- into a separate table later
    
    from bronze.destinations
    ),
-- extra effort to aggregate classification tags (let's review this later, maybe also a separate table?)
class_unnest_1 as (
    select
        identifier,
        unnest(classification, recursive := true)
    from bronze.destinations 
    ),
class_unnest_2 as (
    select 
        identifier,
        name,
        unnest(values) as cls
    from class_unnest_1
    ),
class_final as (
    select
        identifier as id,
        list(cls.title) as classifications
    from class_unnest_2
    where cls.name != 'bysa'
    group by
        identifier
    )
select
    b.*,
    c.classifications
from btbl b
inner join class_final c
    on b.id = c.id;


/* ============== ATTRACTIONS ============== */
create or replace table attractions as
with 
btbl as (
    select
        identifier as id,
        "@type" as record_type,
        name,
        description,
        address, -- explode into separate table later
        geo.latitude as lat,
        geo.longitude as lon,
        image, -- explode later
        containedInPlace.identifier as contained_in_destination_id,
        availableLanguage as available_language,
        event, --clean up later if needed
        isAccessibleForFree as is_accessible_for_free,
        price.minPrice as price, -- there is more, add if needed later
        selfGuided as self_guided,
        reservationRequired as reservation_required,
        email,
        telephone,
        url,
        photo

    from bronze.attractions
    ),
-- extra effort to aggregate classification tags (let's review this later, maybe also a separate table?)
class_unnest_1 as (
    select
        identifier,
        unnest(classification, recursive := true)
    from bronze.attractions 
    ),
class_unnest_2 as (
    select 
        identifier,
        name,
        unnest(values) as cls
    from class_unnest_1
    ),
class_final as (
    select
        identifier as id,
        list(cls.title) as classifications
    from class_unnest_2
    where cls.name != 'bysa'
    group by
        identifier
    )
select
    b.*,
    c.classifications
from btbl b
inner join class_final c
    on b.id = c.id;




/* ============== EXPPORT DB ============== */
EXPORT DATABASE '~/GithubRepos/swisshikefinder/data/db' (
    FORMAT PARQUET,
    COMPRESSION ZSTD,
    ROW_GROUP_SIZE 100_000
);
    
--
--COPY tours to '~/GithubRepos/swisshikefinder/data/db/tours.parquet' (FORMAT PARQUET);
--COPY destinations to '~/GithubRepos/swisshikefinder/data/db/destinations.parquet' (FORMAT PARQUET);
--COPY attractions to '~/GithubRepos/swisshikefinder/data/db/attractions.parquet' (FORMAT PARQUET);


/*random testing*/
create view attractions as select * from read_parquet('/Users/janbenisek/GithubRepos/swisshikefinder/data/db/attractions.parquet');
create view destinations as select * from read_parquet('/Users/janbenisek/GithubRepos/swisshikefinder/data/db/destinations.parquet');
create view tours as  select * from read_parquet('/Users/janbenisek/GithubRepos/swisshikefinder/data/db/tours.parquet');

select * from attractions limit 10;
select * from tours limit 10;
select * from destinations limit 10;

select t.* from tours t where t.itinerary.name = 'Sarnen';
select t.itinerary[1].name from tours t;


describe destinations









