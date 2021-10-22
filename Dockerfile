FROM golang

WORKDIR /go/src/multicast
COPY . .


RUN go get -d -v ./...
RUN go install -v github.com/msalvati1997/b1multicasting/cmd/multicast

CMD ["multicast"]
