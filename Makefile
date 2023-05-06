build:
	./build.sh

run: build
	go-rest-api serve-rest --config dev.config.yml --env=dev

test-server: build
	go-rest-api serve-rest --config test.config.yml --env=dev

run-tests:
	go test ./e2e_test --config=../test.config.yml -ginkgo.v --env=test

serve:
	docker-compose down
	docker-compose up -d
