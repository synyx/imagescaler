.PHONY : go-build docker-build docker-push

go-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o imagescaler .

docker-build:
	docker build -t synyx/imagescaler .

docker-push:
	docker push synyx/imagescaler

publish: go-build docker-build docker-push
