SOURCES=$(shell find . -name '*.go')

test:
	go test ./... -v

bench: testdata/facts.json
	go test ./... -test.run=NONE -test.bench=. -test.benchmem -test.benchtime=5s

benchclean:
	rm testdata/facts.json

graph: testdata/nodes.svg

graphclean:
	rm testdata/nodes.dot

profile: testdata/cpu.profile
	go tool pprof ./qfy.test $<

testdata/cpu.profile: $(SOURCES) testdata/facts.json
	go test . -test.run=NONE -test.bench=. -test.benchmem -test.benchtime=30s -test.cpuprofile=$@

testdata/facts.json: testdata/targets.json
	testdata/generate.rb

testdata/nodes.dot: testdata/targets.json $(SOURCES)
	go run cmd/flood-qfy-graph/main.go -in testdata/targets.json -n 100 > $@

testdata/nodes.svg: testdata/nodes.dot
	dot -Tsvg $< -o $@
