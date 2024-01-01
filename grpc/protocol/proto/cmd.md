hello.proto
protoc --proto_path=. --go_out=. --go_opt=paths=import .\hello.proto
protoc --proto_path=. --go-grpc_out=. --go-grpc_opt=paths=import .\hello.proto

protoc --proto_path=. --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import ./proto/hello.proto

protoc --go_out=.  --go-grpc_out=. ./proto/hello.proto

work.proto
protoc --proto_path=. --go_out=. --go_opt=paths=import .\work.proto
protoc --proto_path=. --go-grpc_out=. --go-grpc_opt=paths=import .\work.proto

product.proto
cd ./proto/
protoc -I ./ -I $GOPATH/src/  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import ./product.proto