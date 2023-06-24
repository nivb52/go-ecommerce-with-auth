
go build ../../../cmd/web/main.go -gcflags=all="-N -l"


go build -o ../../../tmp/api/main.exe ../../../cmd/api/