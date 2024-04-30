/* ============== TRANSFORM ============== 
 * transforms nested data into semi-normalised tables
 * some tables still need further normalisation with:
 *  unnest(additionalDescriptions, recursive := true) 
 * TO-DO: special schema and user?
 * */

create view if not exists tours_raw as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/tours_data.json');
create view if not exists destinations_raw as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/destinations_data.json');
create view if not exists attractions_raw as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/attractions_data.json');


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
create or replace table tours as
with 
btbl as (
    select
        identifier as id,
        "@type" as record_type,
        name,
        abstract,
        touristType[1] as tourist_type, -- mostly 1, sometimes 2 or 3, explode later
        itinerary, -- probably later explode into separate table
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
        specs.barrierFree as barrier_free,
        image[1].url as image_url -- it is alwyas just one (in this dataset at least)
        
    from tours_raw
    ),
-- extra effort to aggregate classification tags (let's review this later, maybe also a separate table?)
class_unnest_1 as (
    select
        identifier,
        unnest(classification, recursive := true)
    from tours_raw 
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
    
    from destinations_raw
    ),
-- extra effort to aggregate classification tags (let's review this later, maybe also a separate table?)
class_unnest_1 as (
    select
        identifier,
        unnest(classification, recursive := true)
    from destinations_raw 
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

    from attractions_raw
    ),
-- extra effort to aggregate classification tags (let's review this later, maybe also a separate table?)
class_unnest_1 as (
    select
        identifier,
        unnest(classification, recursive := true)
    from attractions_raw 
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
    

COPY tours to '~/GithubRepos/swisshikefinder/data/db/tours.parquet' (FORMAT PARQUET);
COPY destinations to '~/GithubRepos/swisshikefinder/data/db/destinations.parquet' (FORMAT PARQUET);
COPY attractions to '~/GithubRepos/swisshikefinder/data/db/attractions.parquet' (FORMAT PARQUET);

