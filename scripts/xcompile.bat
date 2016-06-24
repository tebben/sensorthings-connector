@echo off

if exist ..\bin rd /s /q ..\bin
mkdir ..\bin\win32
mkdir ..\bin\win64
mkdir ..\bin\linux32
mkdir ..\bin\linux64
mkdir ..\bin\linuxARMv7

SET GOOS=windows
SET GOARCH=386
go build -o ..\bin\win32\sensorthings-connector.exe ../src/
echo "Built application for Windows/386"
SET GOARCH=amd64
go build -o ..\bin\win64\sensorthings-connector.exe ../src/
echo "Built application for Windows/amd64"

SET GOOS=linux
SET GOARCH=386
go build -o ..\bin\linux32\sensorthings-connector ../src/
echo "Built application for Linux/386"
SET GOARCH=amd64
go build -o ..\bin\linux64\sensorthings-connector ../src/
echo "Built application for Linux/amd64"
SET GOARCH=arm
SET GOARM=7
go build -o ..\bin\linuxARMv7\sensorthings-connector ../src/
echo "Built application for Linux/ARMv7"
