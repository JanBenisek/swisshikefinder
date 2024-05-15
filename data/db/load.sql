COPY gold.tours FROM './db/gold_tours.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');
COPY gold.tour_classification FROM './db/gold_tour_classification.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');
COPY gold.tour_images FROM './db/gold_tour_images.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');
COPY gold.tour_itinerary FROM './db/gold_tour_itinerary.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');
COPY gold.tour_tourist_types FROM './db/gold_tour_tourist_types.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');
