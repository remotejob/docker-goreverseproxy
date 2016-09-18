all: push

# 0.0 shouldn't clobber any released builds
TAG = 1.5
PREFIX = gcr.io/jntlserv0/docker-goreverseproxy

binary: server.go
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w' -o server

container: binary
	docker build -t $(PREFIX):$(TAG) .

push: container
	gcloud docker push $(PREFIX):$(TAG)

clean:
	docker rmi -f $(PREFIX):$(TAG) || true
