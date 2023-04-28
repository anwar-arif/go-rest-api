build:
	./build.sh

run: build
	go-rest-api serve-rest --config local.config.yaml

test-server: build
	go-rest-api serve-rest --config test.config.yaml

run-tests:
	go test ./e2e_test --config=../test.config.yaml -ginkgo.v

serve:
	docker-compose down
	docker-compose up -d
