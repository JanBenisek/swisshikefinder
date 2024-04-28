# DATA

- The idea is to get all data and convert them to a DB

## Data

- `/destinations` & `destinations/{id}`
  - about destination, links to `myswitzerland.com`
- `/attractions` & `attractions/{id}`
  - about attractions, links to `myswitzerland.com`
- ~~`/tours` & `tours/{id}` & `tours/{id}/geodata`~~: _will get later_
  - has info about tours and hikes, links directly to `schweizmobil.ch`
- ~~`/offers` & `/offers/{id}`~~: _not used_

## Get all data

- `scrape_hikes.py` to obtain all data
- note that we need `extended=true` to get all facets from the bulk endpoint

## Convert to SQL database

- TODO
