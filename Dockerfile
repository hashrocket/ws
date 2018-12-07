FROM golang:1.8

WORKDIR /go/src/github.com/teyushen/ws
COPY . .

CMD ["/bin/bash"]
