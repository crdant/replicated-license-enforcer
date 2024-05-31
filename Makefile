include Makefile.build.mk

.PHONY: test
test: 
	go test -v $(TEST_BUILDFLAGS) ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: build
build: 
	go build ${LDFLAGS} ${GCFLAGS} -v -o bin/enforcer $(BUILDFLAGS) ./cmd/enforcer

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet $(BUILDFLAGS) ./pkg/... ./cmd/...

.PHONY: scan
scan:
	trivy fs \
		--scanners vuln \
		--exit-code=1 \
		--severity="CRITICAL,HIGH,MEDIUM" \
		--ignore-unfixed \
		./
