# Builder image, where we compile/build the executable.
FROM golang:1.9.2 AS builder
WORKDIR /go/src/github.com/timobartels/rest-api
COPY . .
WORKDIR /go/src/github.com/timobartels/rest-api
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
#RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM scratch
LABEL maintainer "Timo Bartels <timo@bartels-online.net>"
COPY --from=builder /go/src/github.com/timobartels/rest-api .
EXPOSE 8080
ENTRYPOINT ["/rest-api"]
