include .env
export

service-run:
	go run main.go

migrate-up:
	migrate -path migrations -database ${CONNECTION} up

migrate-down:
	migrate -path migrations -database ${CONNECTION} down