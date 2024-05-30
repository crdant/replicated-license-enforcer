.PHONY: test
test: 
	go test -v $(TEST_BUILDFLAGS) ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: build
build:
	go build ${LDFLAGS} ${GCFLAGS} -v -o bin/replicated $(BUILDFLAGS) ./cmd/replicated

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet $(BUILDFLAGS) ./pkg/... ./cmd/...

.PHONY: build-ttl.sh
build-ttl.sh:
	docker buildx build .  -t ttl.sh/${USER}/replicated-sdk:24h -f deploy/Dockerfile
	docker push ttl.sh/${USER}/replicated-sdk:24h

	make -C chart build-ttl.sh

.PHONY: mock
mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -source=pkg/store/store_interface.go -destination=pkg/store/mock/mock_store.go

.PHONY: scan
scan:
	trivy fs \
		--scanners vuln \
		--exit-code=1 \
		--severity="CRITICAL,HIGH,MEDIUM" \
		--ignore-unfixed \
		./
