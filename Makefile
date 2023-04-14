build:
	./build.sh

run: build
	go-rest-api serve-rest --config local.config.yaml

serve:
	docker-compose down
	docker-compose up -d
