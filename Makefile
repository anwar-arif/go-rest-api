build:
	./build.sh

run: build
	go-rest-api serve-rest --config example.config.yaml

serve:
	docker-compose down
	docker-compose up -d
