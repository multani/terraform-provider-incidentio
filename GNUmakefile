GOTEST := go test

default: generate testacc

.PHONY: test
test:
	$(GOTEST) ./... -v $(TESTARGS) -timeout 120m

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 $(MAKE) test

.PHONY: generate
generate:
	go generate ./...

