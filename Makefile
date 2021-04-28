TESTABLE_PACKAGES=`go list ./... | grep 'http/controller'`

deps:
	@sh ./dev/deps.sh

run:
	@go run *.go serve

.PHONY: test
test:
	@go test -count=1 -v github.com/guilhermeCoutinho/api-studies/test

test-only:
	@go test -count=1 -v github.com/guilhermeCoutinho/api-studies/test -run ${FILTER}

.PHONY: mocks
mocks:
	@python3 scripts/generate_mocks.py

unit:
	@go test -v ${TESTABLE_PACKAGES} -tags=unit -coverprofile=unit.coverprofile -count=1