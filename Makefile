run: 
	docker compose up --build

newmig:
	migrate create -ext sql -dir db/migrations -seq $(name)

migup:
	migrate -path db/migrations -database "mysql://user:password@tcp(127.0.0.1:3306)/db" -verbose up

migdown:
	migrate -path db/migrations -database "mysql://user:password@tcp(127.0.0.1:3306)/db" -verbose down

bench:
	go test ./api -bench=$(test)

benchmem:
	go test ./api -bench=$(test) -benchmem -memprofile profiles/$(name).out

benchcpu:
	go test ./api -bench=$(test) -cpuprofile profiles/$(name).out

.PHONY: run sqlc migup migdown newmig bench benchmem benchcpu