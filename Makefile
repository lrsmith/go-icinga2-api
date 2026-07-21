default: fmt lint

tag:
	git tag $(shell svu next)
	git push --tags
release: tag

watch:
	watchexec -e go make qa
lint:
	golangci-lint run
fmt:
	gofmt -s -w -e .
docker_start:
	(cd fixtures; docker compose up -d)
	sleep 20
	(cd fixtures; docker compose cp icinga2:/var/lib/icinga2/ca/ca.crt ca.crt)
docker_stop:
	(cd fixtures; docker compose stop; rm -f ca.crt)
test:
	ICINGA2_API_PASSWORD="icingaweb" ICINGA2_API_URL="https://127.0.0.1:5665/v1" ICINGA2_API_USER=icingaweb ICINGA2_INSECURE_SKIP_TLS_VERIFY=true ICINGA2_API_CA_CERT_FILE=$(CURDIR)/fixtures/ca.crt TF_ACC=1 go test -v -cover -timeout 120m ./...
acceptance: docker_start test
qa: lint test
.PHONY: tag release watch lint fmt docker_start docker_stop test acceptance qa
