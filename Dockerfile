FROM golang:1.13 AS builder

ADD . /src
WORKDIR /src

RUN GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /imagescaler .

FROM scratch
COPY --from=builder /imagescaler ./
ENTRYPOINT ["./imagescaler"]
