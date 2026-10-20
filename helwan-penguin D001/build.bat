```bat
@echo off
setlocal
chcp 65001 >nul
title Helwan Penguin - Build

cd /d "%~dp0"

echo.
echo ========================================
echo       Helwan Penguin Build System
echo ========================================
echo.

echo [1/4] Checking Go...
go version
if errorlevel 1 goto error

echo.
echo [2/4] Formatting source files...
gofmt -w cmd internal
if errorlevel 1 goto error

echo.
echo [3/4] Running tests...
go test ./...
if errorlevel 1 goto error

echo.
echo [4/4] Building Windows executable...
if not exist "dist" mkdir "dist"

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

go build -trimpath -o "dist\HelwanPenguin.exe" ./cmd/helwan-penguin
if errorlevel 1 goto error

echo.
echo ========================================
echo BUILD SUCCESSFUL
echo Executable:
echo %CD%\dist\HelwanPenguin.exe
echo ========================================
echo.
pause
exit /b 0

:error
echo.
echo BUILD FAILED. Check the error messages above.
pause
exit /b 1
```