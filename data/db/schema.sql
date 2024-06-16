CREATE SCHEMA bronze;
CREATE SCHEMA gold;
CREATE SCHEMA silver;
CREATE SCHEMA recommendations;

CREATE SEQUENCE rec_id_seq START 1;
CREATE TABLE recommendations.tours(id INTEGER DEFAULT nextval('rec_id_seq'), "title" VARCHAR, "description" VARCHAR, created_at TIMESTAMP);

CREATE TABLE gold.tours(ID UUID, record_type VARCHAR, "name" VARCHAR, abstract VARCHAR, part_of_trip_id UUID, url_mysw VARCHAR, lat DOUBLE, lon DOUBLE, elevation BIGINT, logo_url VARCHAR, stage VARCHAR, requirements_technical VARCHAR, requirements_endurance VARCHAR, route_category VARCHAR, url_swmo VARCHAR, season VARCHAR, distance BIGINT, duration BIGINT, duration_reverse_direction BIGINT, ascent BIGINT, descent BIGINT, barrier_free BOOLEAN);
CREATE TABLE gold.tour_classification(ID UUID, context VARCHAR, "name" VARCHAR, classification VARCHAR);
CREATE TABLE gold.tour_images(ID UUID, "type" VARCHAR, keywords VARCHAR, publisher VARCHAR, encoding_format VARCHAR, url VARCHAR, "name" VARCHAR, copyright_holder VARCHAR, width BIGINT, height BIGINT);
CREATE TABLE gold.tour_itinerary(ID UUID, "type" VARCHAR, itinerary_id UUID, "name" VARCHAR, url VARCHAR);
CREATE TABLE gold.tour_tourist_types(ID UUID, tourist_type VARCHAR);

CREATE VIEW bronze.tours AS SELECT * FROM read_json('./raw_data/tours_data.json');




