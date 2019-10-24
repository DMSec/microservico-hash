FROM golang:1.11-stretch

ADD . $GOROOT/src/microservicos-hash

RUN go get -u google.golang.org/grpc \
 && go get -u github.com/golang/protobuf/proto \
 && go get -u github.com/golang/protobuf/protoc-gen-go \
 && go get -u github.com/go-sql-driver/mysql


WORKDIR $GOROOT/src/microservicos-hash/listagem

RUN ls $GOROOT/src/microservicos-hash/listagem
RUN go env
RUN go install

CMD ["go", "run", "main.go"]
