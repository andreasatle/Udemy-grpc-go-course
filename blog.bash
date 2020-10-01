#/bin/bash
protoc blog/pb/blog.proto --go_out=plugins=grpc:.
