# SDCC
SDCC project
##How to build
In order to compile this program you need Go installed.
###Executable
The executable can be built with:
```sh
make executable
```
###Proto
For this section you need protoc/protobuf installed.

Install the protocol compiler plugins for Go using the following commands:
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```
Update your PATH so that the protoc compiler can find the plugins:
```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```
You can now generate protobuf related code:
```sh
make proto
```
The following command generates both protobuf code and server executable:
```sh
make all
```
