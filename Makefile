deps:
	@sh ./dev/deps.sh
	@make migrate

migrate:
	@pushd migrations; make migrate; popd

run-server:
	@go run *.go serve