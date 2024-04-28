--create view tours_v as select * from read_json('/Users/janbenisek/GithubRepos/swisshikefinder/data/results/tours_data.json');

select
    identifier as id,
    "@type" as record_type,
    name,
    abstract,
    touristType[1] as tourist_type, -- assuming there is always just 1 type in the list
    itinerary,
    partOfTrip.identifier as part_of_trip_id,
    url as url_mysw,
    geo.latitude as lat,
    geo.longitude as lon,
    geo.elevation as elevation,
    -- switzerland mobility info
    switzerlandMobility.logo as logo_url,
    switzerlandMobility.stage as stage,
    switzerlandMobility.requirements as requirements,
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
    image[1].url as image_url, -- assuming there's just 1 image there, otherwise take just the first
    classification
    
from tours_v
limit 100;

--describe tours_v;

with btbl as (
    select
        identifier,
        unnest(classification, recursive := true)
    from tours_v 
    where identifier = 'ffcf3a7e-c62b-48c3-8e43-67c80e7fbc56' 
    limit 5
    ),
second_level as (
    select 
        identifier,
        name,
        unnest(values) as cls
    from btbl

),
extract_classification as (
    select 
        identifier,
        name,
        cls.*
        
    from second_level
)
--select * from extract_classification

-- wide format
--select 
--    * 
--    from (
--        pivot extract_classification 
--            on name 
--            using any_value(title) 
--            group by identifier
-- )
-- agg into list
select
    identifier,
    name,
    list(title) as classifications
from extract_classification
group by
    identifier,
    name
;



