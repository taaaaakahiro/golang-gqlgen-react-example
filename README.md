# Golang Example API

## How to run

### Setup environment

1. Make sure that you have `direnv` installed to configure local environment variables. Please look at the [direnv github](https://github.com/direnv/direnv#install) for installation.
Copy the `.envrc.example` file and set the environment variables, then enable `direnv`.

```console
$ cp .envrc.sample .envrc
$ vi .envrc
$ direnv allow
```

2. Make & edit `.env` in API directories.
```console
$ cd api_xxx
$ cp .env.sample .env
$ vi .env
```

3. And then, start servers and databases by running the command below:

```console
$ docker-compose up
```

---
# ALL Docker Containers

## Mongo
### Mongo DB (*comment out*)
- port: 27017

### API (*comment out*)
- port: 8080
- url: http://0.0.0.0:8080/

### React without any lib (*comment out*)
- port: 3001
- url: http://0.0.0.0:3031/

### React with Relay (*un-used*)
- [Relay](https://relay.dev/)
- port: 3002
- url: http://0.0.0.0:3002/

## MySQL8

### MySQL DB
- port: 3306
- X API port(33060) is closed. (golang's driver unsupported X API)

### API
- port: 8081
- url: http://0.0.0.0:8081/

### React without any lib
- port: 3003
- url: http://0.0.0.0:3003/

---
# Run API Test

- See README.md in `/api_xxx`.
