Setup in order to run gprc with go
==================================

brew install protobuf

go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
go install google.golang.org/grpc
go install github.com/golang/protobuf/protoc-gen-go

GRPC_PATH=${GOPATH}/src/github.com/andreasatle/Udemy/grpc-go-course
GRPC_GREET_PATH=${GRPC_PATH}/greet
GRPC_CALC_PATH=${GRPC_PATH}/calculator 
GRPC_BLOG_PATH=${GRPC_PATH}/blog

GRPC_XXX_PB_PATH=${GRPC_PATH}/XXX/XXXpb
GRPC_XXX_SERVER_PATH=${GRPC_PATH}/XXX/XXX_server
GRPC_XXX_CLIENT_PATH=${GRPC_PATH}/XXX/XXX_client

XXX = greet or calculator

To use reflection with evans CLI:
go get -u github.com/ktr0731/evans

To start mongoDB:
brew tap mongodb/brew
brew install mongodb-community@4.4
brew services start mongodb/brew/mongodb-community 

go get -u go.mongodb.org/mongo-driver/mongo
[go get -u github.com/mongodb/mongo-go-driver/mongo]
mongod --dbpath=/Users/andreasatle/GoLang/src/github.com/andreasatle/Udemy/grpc-go-course/data/db
