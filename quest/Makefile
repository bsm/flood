SOURCES=$(shell find . -name '*.go')

test:
	go test ./... -v

bench:
	go test ./... -run=NONE -bench=. -benchmem -benchtime=5s

profile-mem: testdata/mem.profile
	go tool pprof --alloc_objects ./quest.test $<

testdata/mem.svg: testdata/mem.profile
	go tool pprof --alloc_objects -cum -focus=Match -svg -output=$@ ./quest.test $<

testdata/mem.profile: $(SOURCES)
	go test . -run=NONE -bench=. -benchmem -benchtime=30s -memprofile=$@
