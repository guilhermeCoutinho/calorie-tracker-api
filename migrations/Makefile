ADDRESS := localhost:9000
PASSWORD := ""
DATABASE_NAME := "calorie-tracker"
DATABASE_USER := "calorie-tracker-user"

migrate:
	@go run *.go migrate -address $(ADDRESS) -pass $(PASSWORD) -database $(DATABASE_NAME) -user $(DATABASE_USER)

rollback:
	@go run *.go rollback -address $(ADDRESS) -pass $(PASSWORD) -database $(DATABASE_NAME) -user $(DATABASE_USER)