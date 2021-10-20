FROM golang

WORKDIR /go/src/multicast
COPY . .


RUN go get -d -v ./...
RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee
RUN go get github.com/prometheus/client_golang/prometheus@v1.7.0
RUN go install -v github.com/msalvati1997/b1multicasting/cmd/multicast

CMD ["multicast"]
