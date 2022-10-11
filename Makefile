.PHONY: test

test:
	go test ./test

clean:
	rm -rf ./bin ./vendor Gopkg.lock

check: test clean
	go vet .