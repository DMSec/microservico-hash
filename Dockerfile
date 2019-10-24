FROM golang:1.11-stretch

ADD . $GOROOT/src/microservicos-hash

RUN go get google.golang.org/grpc \
 && go get github.com/golang/protobuf/proto \
 && go get github.com/golang/protobuf/protoc-gen-go \
 && go get github.com/go-sql-driver/mysql \
 && go get github.com/DMSec/microservico-hash/listagem


WORKDIR $GOROOT/src/microservicos-hash/listagem

RUN ls $GOROOT/src/microservicos-hash/listagem
RUN go env

WORKDIR $GOROOT/src/microservicos-hash/listagem

CMD ["go", "run", "main.go"]
