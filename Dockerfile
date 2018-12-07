FROM golang:1.8

WORKDIR /go/src/github.com/hashrocket/ws
COPY . .


RUN go get -d -v github.com/hashrocket/ws
RUN go install -v github.com/hashrocket/ws

CMD ["/bin/bash"]
