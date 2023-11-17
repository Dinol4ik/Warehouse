up:
	 docker-compose up

migrate:
	goose postgres "host=localhost user=user password=postgres dbname=lamoda sslmode=disable" up

rollback:
	goose postgres "host=localhost user=user password=postgres dbname=lamoda sslmode=disable" down

testGo:
	go test .\internal\controller\. -v