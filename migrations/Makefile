ADDRESS := localhost:9000
PASSWORD := ""
DATABASE_NAME := racing-clash
DATABASE_USER := racing-clash-user
migrate:
	@cd migrations && go run *.go  -address $(ADDRESS) -pass $(PASSWORD) -database $(DATABASE_NAME) -user $(DATABASE_USER)
migrate-rollback-all:
	@cd migrations && go run *.go  -address $(ADDRESS) -pass $(PASSWORD) -database $(DATABASE_NAME) -user $(DATABASE_USER) reset
migrate-rollback-last:
	@cd migrations && go run *.go  -address $(ADDRESS) -pass $(PASSWORD) -database $(DATABASE_NAME) -user $(DATABASE_USER) down
migrate-version:
	@cd migrations && go run *.go  -address $(ADDRESS) -pass $(PASSWORD) -database $(DATABASE_NAME) -user $(DATABASE_USER) version