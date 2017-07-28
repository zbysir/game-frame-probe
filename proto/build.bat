protoc -I=. -I=%GOPATH%\src --gogoslick_out=plugins=grpc:./pbgo ./*.proto 
pause