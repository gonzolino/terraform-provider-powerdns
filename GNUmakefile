GOSWAGGER_IMAGE=quay.io/goswagger/swagger
GOSWAGGER_VERSION=v0.28.0
SWAGGERCMD=docker run --rm -v $(HOME):$(HOME) -w $(CURDIR) $(GOSWAGGER_IMAGE):$(GOSWAGGER_VERSION)
SWAGGER_SPEC_FILE=swagger.yaml

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

$(SWAGGER_SPEC_FILE):
	# Use swagger spec from openHAB
	curl -o $(SWAGGER_SPEC_FILE) https://raw.githubusercontent.com/PowerDNS/pdns/master/docs/http-api/swagger/authoritative-api-swagger.yaml

.PHONY: generate
generate: $(SWAGGER_SPEC_FILE)
	go generate
	$(SWAGGERCMD) generate client -f $(SWAGGER_SPEC_FILE) -c internal/powerdns/client -m internal/powerdns/models
