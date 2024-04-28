# DATA

- The idea is to get all data and convert them to a DB

## Get all data

- `scrape_hikes.py`
- First I get IDs of all objects from each API
- Then get all elements from `<API>/{id}`
- this is because there are rate limits and only the `<API>/{id}` endpoint has the full data

## Convert to SQL database

- TODO
