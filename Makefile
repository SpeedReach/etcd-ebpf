all: docker

binary:
	go build -o bin/main internal/cmd/main.go

docker: binary
	docker build -t github.com/speedreach/ebpf-etcd .