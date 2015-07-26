default: test

test:
	go test ./...

bench:
	go test ./... -test.run=NONE -test.bench=. -test.benchmem
