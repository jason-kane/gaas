GCO_ENABLED=0
GOARCH=amd64
GOOS=linux

build:
	go build -x -a -tags netgo -ldflags '-w' -o gaas .
	# -installsuffix cgo
	docker build -t gaas .

run:
	docker run --publish 8080:8080 gaas:latest

shell:
	docker run -ti gaas:latest sh