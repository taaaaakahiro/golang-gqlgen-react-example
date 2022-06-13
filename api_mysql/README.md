# Golang Example API

## How to run

### Setup environment

1. Make sure that you have `direnv` installed to configure local environment variables. Please look at the [direnv github](https://github.com/direnv/direnv#install) for installation.
Copy the `.envrc.example` file and set the environment variables, then enable `direnv`.

```console
$ cd api_mysql
$ cp .envrc.sample .envrc
$ vi .envrc
$ direnv allow
```

2. Make & edit `.env` in API directories.
```console
$ cd api_mysql
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
[shema.graphqls](/api_mysql/pkg/graph/schema.graphqls)
- query
  - users
    - get all users
  - messages
    - get messages by user
  - allMessages
    - get all messages
- mutation
  - createMessage
    - insert a message by user

### How to edit GraphQL schema

- [gqlgen](https://gqlgen.com/getting-started/)
- ref: [https://tech.layerx.co.jp/entry/2021/10/22/171242](https://tech.layerx.co.jp/entry/2021/10/22/171242)

1. Edit schema in `/api_mysql/pkg/graph/schema.graphqls`.
2. Run command.
   ```
   $ cd api_mysql
   $ go run github.com/99designs/gqlgen generate
   ```
3. Edit `/api_mysql/pkg/graph/schema.resolvers.go`.

---
# Docker Containers

## MySQL DB
- port: 3306
- X API port(33060) is closed. (golang's driver unsupported X API)
- MySQL 8

## API
- port: 8081
- url: http://0.0.0.0:8081/
- golang
- gqlgen
- dataloader

---
# Run API Test

1. Go to API test directory
```console
$ cd api_mysql/test_fixtures/
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

3. Run database docker container.
```console
$ make up
```

4. Run go test.
```console
$ make test
```

5. Refresh go test cache & DB.
```console
$ make refresh
```

6. Finish test.
```console
$ make stop
```

7. Down database container.
```console
$ make down
```