generate:
	go generate ./gen

test:
	go test ./tests -run TestArenaEqualsStdJSON -v

bench:
	go test ./tests -bench BenchmarkArenaPipeline -benchmem