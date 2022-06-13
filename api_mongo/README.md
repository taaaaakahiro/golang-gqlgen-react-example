# Golang Example API

## How to run

### Setup environment

1. Make sure that you have `direnv` installed to configure local environment variables. Please look at the [direnv github](https://github.com/direnv/direnv#install) for installation.
Copy the `.envrc.example` file and set the environment variables, then enable `direnv`.

```console
$ cd api_mongo
$ cp .envrc.sample .envrc
$ vi .envrc
$ direnv allow
```

2. Make & edit `.env` in API directories.
```console
$ cd api_mongo
$ cp .env.sample .env
$ vi .env
```

3. And then, start servers and databases by running the command below:

```console
$ cd ../
$ docker-compose up
```

## API Interfaces

| Func               | Method | Path                   | Description                            |
|--------------------|--------|------------------------|----------------------------------------|
| GetUsers           | GET    | /user/list             | get list by all users                  |
| GetMessages        | GET    | /message/list/:user_id | get all messages for specified user id |
| GetVersion         | GET    | /version               | (common) get API version               |
| healthCheckHandler | GET    | /healthz               | (common) return HTTP Status 200        |
| Query              | POST   | /query                 | GraphQL endpoint                       |
| -                  | GET    | /gql                   | GraphQL playground                     |

---
## GraphQL

### Schema definitions
[shema.graphqls](/api_mongo/pkg/graph/schema.graphqls)
- query
   - users
      - get all users
   - messages
      - get messages by user
- mutation
   - createMessage
      - insert a message by user

### How to edit GraphQL schema

- [gqlgen](https://gqlgen.com/getting-started/)
- ref: [https://tech.layerx.co.jp/entry/2021/10/22/171242](https://tech.layerx.co.jp/entry/2021/10/22/171242)

1. Edit schema in `api_mongo/pkg/graph/schema.graphqls`.
2. Run command.
   ```
   $ cd api_mongo
   $ go run github.com/99designs/gqlgen generate
   ```
3. Edit `api_mongo/pkg/graph/schema.resolvers.go`.

---
# Docker Containers

## Mongo DB
- port: 27017

## API
- port: 8080
- url: http://0.0.0.0:8080/
- golang
- gqlgen

---
# Run API Test

1. Go to API test directory
```console
$ cd api_mongo/test_fixtures/
```

2. (first time) Prepare test.
   1. Allow direnv in API test directories.
   ```console
   $ direnv allow
   ```
   2. Make & edit `.env` in API directories.
   ```console
   $ cp .env.sample .env
   $ vi .env
   ```

3. Run database docker
```console
$ make up
```

4. Run go test
```console
$ make test
```

5. Finish test
```console
$ make stop
```
6. Down docker
```console
$ make down
```
