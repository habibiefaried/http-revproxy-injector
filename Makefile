test:
	go fmt ./...
	go build
	docker-compose build
	docker-compose up -d
	docker exec unittest bash -c "GITHUB_RUN_ID=1 go test -p 1 -v ./..."