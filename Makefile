.PHONY : go-build docker-build docker-push

go-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o imagescaler .

docker-build:
	docker build -t rjayasinghe/imagescaler .

docker-push:
	docker push rjayasinghe/imagescaler

publish: go-build docker-build docker-push
