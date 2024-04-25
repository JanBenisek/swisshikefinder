# Swiss Hike Finder

Playground to learn Go and find awesome hikes.

## Next Steps

- slightly prettier UI
- clean up code and review structure
- add logging
- buy domain and switch
- deployment - droplet with docker compose?

## Go notes

- `:=` declares and inits a variable and it infers the type
- `=` assigns value to an existing variable
- `&` creates a new pointer to a new instance, useful to pass a pointer instead of an instance
  - benefit is mutability - if I need to mutate variable elsewhere, pointer helps
  - go passes arguments to functions as values (so changes would not be visible elsewhere)
  - whenever we pass an instance of a struct, go makes a copy, which might be inefficient
- `composite literals` are made with `{...}`
  - to create and initialise values for composite types, like arrays, structs, maps etc..
  - when go creates new instances of these, they are initialised as pointers

## Readings

- https://freshman.tech/web-development-with-go/
- https://stackoverflow.com/questions/47270595/how-to-parse-json-string-to-struct
- https://docs.docker.com/language/golang/build-images/
- https://dev.to/willvelida/pushing-container-images-to-github-container-registry-with-github-actions-1m6b
- https://dev.to/francescoxx/build-a-crud-rest-api-in-go-using-mux-postgres-docker-and-docker-compose-2a75