.default: run

run: build
	./amfi-fetch 2> ./err

build: format main.go
	go build -o amfi-fetch main.go

format: main.go amfi/amfi.go
	goimports -l -w ./
