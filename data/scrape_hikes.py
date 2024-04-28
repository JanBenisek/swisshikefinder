# %%
import httpx
import os
from dotenv import load_dotenv
import json
from typing import Dict, Tuple
import time

load_dotenv()

HIKE_API_KEY = os.environ.get("HIKE_API_KEY")

'''
- first get IDs of every dest/att/tour (should be within 1k/day)
- then better control where to start/end

- get all from destination{id}
- get all from attraction{id}
- get all from tour{id} and geodata
'''


class ScrapeHikes:
    # TODO: replace print by logger
    # TODO: try/except for better error handling
    def __init__(self, api_key: str):
        self.api_key = api_key

    def write(self, data, file_name: str):
        '''write file to destination (overwrites it!)'''
        write_path = f"./results/{file_name}.json"
        print(f"Saved to {write_path}")
        with open(write_path, "w") as fp:
            json.dump(data, fp)

    # def read(self, path: str):
    #     with open("./results/destination_ids.json", "w") as fp:
    #         result = json.load(fp)
    #     return result

    def call_api(self, api: str, page: int=0) -> Tuple[Dict, str]:
        ''' get response from API'''
        response = httpx.get(
            f"https://opendata.myswitzerland.io/v1/{api}", 
            headers={
                "accept": "application/json",
                "x-api-key": self.api_key
            }, 
            params={
                "facets": "*",
                "lang": "en",
                "hitsPerPage": 50,
                "striphtml": "true",
                "page": page
            }
        )
        return response.json(), response.status_code

    def get_all_ids(self, api: str, write_to: str):
        '''get IDs of all object in given API'''
        self.result_ids = []
        current_page = 0
        not_last_page = True

        while not_last_page:
            print(f"Processing {api}, page: {current_page}")
            
            result, status = self.call_api(api=api, page=current_page)
            total_pages = result['meta']['page']['totalPages']
            
            print(f" ... Request status: {status}, total pages: {total_pages}")

            # add all ids
            self.result_ids.extend([record['identifier'] for record in result['data']])
            
            current_page += 1
            # 0-based indexing, (83 pages, first with 0 last is 82)
            if (current_page) == total_pages:
                not_last_page = False

            # api allows 1 call per second, let's be nice
            time.sleep(1.5)

        self.write(data=self.result_ids, file_name=write_to)
        print('Finished')

# %%
SH = ScrapeHikes(api_key=HIKE_API_KEY)

# %%
# destinations has 4104 elements, 83 pages with 50 elements per page
SH.get_all_ids(api='destinations', write_to='destination_ids')

#%%
# attractions has 3688 elements, 74 pages with 50 elements per page
SH.get_all_ids(api='attractions', write_to='attraction_ids')

# %%
# tours has 2493 elements, 50 pages with 50 elements per page
SH.get_all_ids(api='tours', write_to='tour_ids')
