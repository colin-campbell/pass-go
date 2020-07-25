.PHONY: build clean docker run push
TAG := digitalist-se/pass-go
clean:
	rm -f pass-go pkged.go

build: clean
	go mod download
	go mod verify
	go generate
	# Cross-compile for linux
	GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-s -w -extldflags "-static"' .

docker: build
	docker build --progress=plain --no-cache -t ${TAG} .

run: docker
	docker run --rm -it -p 5000:5000 ${TAG}