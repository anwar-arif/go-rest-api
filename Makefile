build:
	./build.sh

run: build
	go-rest-api serve-rest --config local.config.yml --env=local

test-server: build
	go-rest-api serve-rest --config test.config.yml --env=test

run-tests:
	go test ./e2e_test --config=../test.config.yml -ginkgo.v --env=test

serve-container:
	docker-compose -f docker-compose-dev.yml down
	docker-compose -f docker-compose-dev.yml up -d
