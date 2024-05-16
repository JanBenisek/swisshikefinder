# DATA

- The idea is to get all data and convert them to a DB
- UPDATE: I do not want to pay for hosted postgres, I'll do in-memory duckDB (for now)

## Data

- `/destinations` & `destinations/{id}`
  - about destination, links to `myswitzerland.com`
- `/attractions` & `attractions/{id}`
  - about attractions, links to `myswitzerland.com`
- `/tours` & `tours/{id}` & ~~`tours/{id}/geodata`~~: _will get later_
  - has info about tours and hikes, links directly to `schweizmobil.ch`
- ~~`/offers` & `/offers/{id}`~~: _not used_

## Get all data

- `scrape_hikes.py` to obtain all data
- note that we need `expand=true` to get all facets from the bulk endpoint

## Installation

- On Mac OS, run: `brew install poetry pyenv pyenv-virtualenv`.
- Add to `~/.zshrc`:

```bash
eval "$(pyenv init --path)"
eval "$(pyenv virtualenv-init -)"
```

- install python version (order matters):
- then `pyenv activate swisshikefinder-3.11.8`

```bash
export PYTHON_VERSION=3.11.8
export PROJECT_NAME=swisshikefinder

pyenv install --skip-existing $PYTHON_VERSION
pyenv virtualenv $PYTHON_VERSION $PROJECT_NAME-$PYTHON_VERSION || true
pyenv local $PROJECT_NAME-$PYTHON_VERSION

poetry install
```

- formatting

```shell
ruff check scrape_hikes.py
ruff format scrape_hikes.py
```
