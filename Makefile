.PHONY: unittests

unit:
	go test -v -tags="unit" github.com/aladdinid/fabric-devkit/internal/config
	go test -v -tags="unit" github.com/aladdinid/fabric-devkit/internal/configtx
	go test -v -tags="unit" github.com/aladdinid/fabric-devkit/internal/crypto
	go test -v -tags="unit" github.com/aladdinid/fabric-devkit/internal/network

smoke:
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/internal/docker
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/internal/config
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/internal/configtx
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/internal/crypto
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/internal/network

test:
	make unit
	make smoke