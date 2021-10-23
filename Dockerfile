FROM golang

WORKDIR /go/src/multicast
COPY . .

RUN go mod download
RUN go clean
RUN go mod tidy
RUN go get -d -v ./... \
    && go install -v github.com/msalvati1997/b1multicasting/cmd/multicast

CMD ["multicast"]
