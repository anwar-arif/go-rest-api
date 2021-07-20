build:
	./build.sh

run: build
	go run main.go serve-rest --config ./example.config.yaml

serve:
	docker-compose down
	docker-compose up -d
