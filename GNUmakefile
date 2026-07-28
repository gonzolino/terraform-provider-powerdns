OPENAPI_SPEC_FILE=openapi.yaml

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

$(OPENAPI_SPEC_FILE):
	curl -fsSL -o $(OPENAPI_SPEC_FILE).tmp https://raw.githubusercontent.com/PowerDNS/pdns/auth-4.9.4/docs/http-api/openapi/authoritative-api-openapi.yaml && mv $(OPENAPI_SPEC_FILE).tmp $(OPENAPI_SPEC_FILE)

.PHONY: generate
generate: $(OPENAPI_SPEC_FILE)
	go generate ./...
