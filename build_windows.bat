set GOARCH=amd64
set GOOS=windows
set GOPATH=%GOPATH%;%cd%;

go build -o bin/download_windows.exe ./src/download.go
go build -o bin/parser_windows.exe ./src/parser.go
