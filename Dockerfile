FROM golang:1.11 AS builder

WORKDIR /src

RUN git clone  https://github.com/rjayasinghe/imagescaler

WORKDIR /src/imagescaler
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /imagescaler .

FROM scratch
COPY --from=builder /imagescaler ./
ENTRYPOINT ["./imagescaler"]
