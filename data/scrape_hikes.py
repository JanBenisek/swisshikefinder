# %%
import httpx
import os
from dotenv import load_dotenv

load_dotenv()

HIKE_API_KEY = os.environ.get("HIKE_API_KEY")

'''
- first get IDs of every dest/att/tour (should be within 1k/day)
- then better control where to start/end

- get all from destination{id}
- get all from attraction{id}
- get all from tour{id} and geodata
'''

url = "https://opendata.myswitzerland.io/v1/destinations"
headers = {
    "accept": "application/json",
    "x-api-key": HIKE_API_KEY
}
params = {
    "facets": "*",
    "lang": "en",
    "hitsPerPage": 50,
    "striphtml": "true"
}

response = httpx.get(url, headers=headers, params=params)

out = response.json()

current_page_num = out['meta']['page']['number']
total_pages = out['meta']['page']['totalPages'] # has 83 pages, but we start from 0 and end at 82
current_data = out['data']
next_link_to_data = out['link']['next']

d = out['data']



# %%
with httpx.Client() as client:
    response = client.get(url, headers=headers, params=params)

    if response.status_code == 200:
        # Handle successful response
        print("HTTPX call successful!")
        print(response.json())  # Print the JSON response content
    else:
        # Handle error
        print("HTTPX call failed with status code:", response.status_code)

out = response.json()
d = out['data']

