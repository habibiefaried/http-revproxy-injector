test:
	go fmt ./...
	go build
	docker-compose build
	docker-compose up -d
	docker exec unittest go test -v ./...
	docker-compose down