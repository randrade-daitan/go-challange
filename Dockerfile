FROM golang:alpine
RUN apk add git

RUN mkdir /challenge
ADD . /challenge
WORKDIR /challenge/cmd/grpc

RUN go mod download
RUN go build -o app .

CMD ["/challenge/cmd/grpc/app"]