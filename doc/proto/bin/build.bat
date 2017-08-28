protoc -I=../ -I=%GOPATH%/src --gogoslick_out=plugins=grpc:../../../common/pbgo ../*.proto
pause