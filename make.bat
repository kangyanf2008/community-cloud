::编译linux版本
cd src
set GOOS=linux
set GOARCH=amd64
set GOHOSTOS=linux
go.exe build  -o ../bin/community-cloud main/main.go
cd ../