# Go rest api

## Run on local
### Build application binary
```bash
$ ./build.sh
or
$ make build
```


### Run application binary (run server)
```bash
$ go-rest-api serve-rest
or
$ make run
```

## Run on Container
```bash
$ docker-compose up --build
or
$ make serve
```

## Testing
### unit test
```bash
$ go test $(go list ./... | grep -v go-rest-api/e2e_test) -v
```
### end to end test
```bash
$ cd e2e_test
$ go test --config=../test.config.yml -ginkgo.v
```

## Folder structure

* api folder contains rest controllers, middlewares
* cmd folder contains application's base like main files
* config folder contains app configuration files
* e2e_test folder contains end-to-end testing suits
* infra contains drivers like db, messaging, cache etc
* repo folder contains database code
* model folder contains model
* service folder contains application service

### flow
> cmd -> api -> service -> repo, models, cache, messaging


### Example APIS

### Health 

Method: `GET`
URL: `http://{base_url}:{system_server_port}/system/v1/health/api`

Response:
```status_code: 200```
```json=
{
    "data": "ok"
}
```

### API list
Method: `GET` URL: `http://localhost:8000`
Response: ```status_code: 200```
```json
{
  "message": "success",
  "data": {
    "method": "GET",
    "service_name": "Go rest api"
  }
}

```