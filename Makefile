deps:
	@sh ./dev/deps.sh

run:
	@go run *.go serve

.PHONY: test
test:
	@go test -count=1 -v github.com/guilhermeCoutinho/api-studies/test

test-only:
	@go test -count=1 -v github.com/guilhermeCoutinho/api-studies/test -run ${FILTER}