PKG:=$(shell glide nv)

default: test

test:
	go test $(PKG)

bench:
	go test $(PKG) -test.run=NONE -test.bench=. -test.benchmem
