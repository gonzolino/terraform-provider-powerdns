OPENAPI_SPEC_FILE=openapi.yaml

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

$(OPENAPI_SPEC_FILE):
	curl -o $(OPENAPI_SPEC_FILE) https://raw.githubusercontent.com/PowerDNS/pdns/master/docs/http-api/openapi/authoritative-api-openapi.yaml

.PHONY: generate
generate: $(OPENAPI_SPEC_FILE)
	go generate ./...
