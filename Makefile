run-pg:
	docker run --name my_postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -e POSTGRES_HOST_AUTH_METHOD=trust -d -p 5432:5432 postgres || true
build:
	cd cmd/gophermart && rm gophermart-binary || true && go build -o gophermart-binary *.go
run-gophermart:
	make run-pg && make build && ./cmd/gophermart/gophermart-binary -l=DEBUG -a=localhost:8080