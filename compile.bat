@echo off
go mod tidy
go build -o app.exe
app.exe