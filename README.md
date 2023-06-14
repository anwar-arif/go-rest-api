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
$ go-rest-api serve-rest -c local.config.yml --env=local
or
$ make run
```
### unit test
```bash
$ go test $(go list ./... | grep -v go-rest-api/e2e_test) -v
```
### end to end test
- run test server
```bash
$ ./build.sh && go-rest-api serve-rest --config test.config.yml --env=test
or 
$ make test-server
```
- in seperate terminal window run the test suites
```bash
$ go test ./e2e_test --config=../test.config.yml -ginkgo.v --env=test
or 
$ make run-tests 
```


## Run on Container
```bash
$ docker-compose -f docker-compose-dev.yml up --build
or
$ make serve-container
```

## Run on local kubernetes cluster
Steps are described in `k8s.README.md` file

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