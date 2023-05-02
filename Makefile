build:
	./build.sh

run: build
	go-rest-api serve-rest --config local.config.yml

test-server: build
	go-rest-api serve-rest --config test.config.yml

run-tests:
	go test ./e2e_test --config=../test.config.yml -ginkgo.v

serve:
	docker-compose down
	docker-compose up -d
