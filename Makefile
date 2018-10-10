.PHONY: unittests

unittests:
	go test -v github.com/aladdinid/fabric-devkit/internal/config
	go test -v github.com/aladdinid/fabric-devkit/internal/configtx
	go test -v github.com/aladdinid/fabric-devkit/internal/crypto
	go test -v github.com/aladdinid/fabric-devkit/internal/network

smoketests:
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/internal/docker