run-pg:
	docker run --name my_postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -e POSTGRES_HOST_AUTH_METHOD=trust -d -p 5432:5432 postgres || true
build:
	cd cmd/gophermart && rm gophermart-binary || true && go build -o gophermart-binary *.go
run-accrual:
	./cmd/accrual/accrual_darwin_arm64
run-gophermart:
	make run-pg && make build && ./cmd/gophermart/gophermart-binary -l=INFO -a=localhost:8080 -r=localhost:8081
test:
	make build && ./gophermarttest-darwin-arm64 \
		-test.v -test.run=^TestGophermart$$ \
		-gophermart-binary-path=cmd/gophermart/gophermart-binary \
		-gophermart-host=localhost \
		-gophermart-port=8081 \
		-gophermart-database-uri="postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable" \
		-accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
		-accrual-host=localhost \
		-accrual-port=8080 \
		-accrual-database-uri="postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
get-coverage:
	go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./... \
    && go tool cover -func coverage.out | grep total | awk '{print $3}'