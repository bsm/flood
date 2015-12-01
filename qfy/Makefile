SOURCES=$(shell find . -name '*.go')

test:
	go test ./... -v

bench: testdata/facts.json
	go test ./... -run=NONE -bench=. -benchmem -benchtime=5s

benchclean:
	rm testdata/facts.json

profile-cpu: testdata/cpu.profile
	go tool pprof ./qfy.test $<

profile-mem: testdata/mem.profile
	go tool pprof --alloc_objects ./qfy.test $<

testdata/cpu.svg: testdata/cpu.profile
	go tool pprof -cum -focus=Select -svg -output=$@ ./qfy.test $<

testdata/mem.svg: testdata/mem.profile
	go tool pprof --alloc_objects -cum -focus=Select -svg -output=$@ ./qfy.test $<

testdata/cpu.profile: $(SOURCES) testdata/facts.json
	go test . -run=NONE -bench=. -benchmem -benchtime=30s -cpuprofile=$@

testdata/mem.profile: $(SOURCES) testdata/facts.json
	go test . -run=NONE -bench=. -benchmem -benchtime=30s -memprofile=$@

testdata/facts.json: testdata/targets.json
	testdata/generate.rb
