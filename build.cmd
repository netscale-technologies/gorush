@ECHO OFF
FOR /F "tokens=* USEBACKQ" %%F IN (`git describe --tags --always`) DO (
SET VERSION=%%F
)
go build -o bin/gorush.exe -v -ldflags 'extldflags=-s,-w'