FROM golang:1.11 AS builder

RUN go get github.com/rjayasinghe/imagescaler

WORKDIR $GOPATH/src/github.com/rjayasinghe/imagescaler
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /imagescaler .

FROM scratch
COPY --from=builder /imagescaler ./
ENTRYPOINT ["./imagescaler"]
