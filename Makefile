.PHONY: unittests

unit:
	go test -v -tags="unit" github.com/aladdinid/fabric-devkit/maejor/svc

smoke:
	go test -v -tags="smoke" github.com/aladdinid/fabric-devkit/maejor/svc

test:
	make unit
	make smoke