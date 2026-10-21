
@echo off
setlocal
chcp 65001 >nul
title Helwan Penguin - Build

cd /d "%~dp0"

echo.
echo ========================================
echo        Helwan Penguin Build System
echo        9 Worlds - 45 Levels
echo ========================================
echo.

echo [1/5] Checking Go...
go version
if errorlevel 1 goto error

echo.
echo [2/5] Formatting source files...
gofmt -w cmd internal tests
if errorlevel 1 goto error

echo.
echo [3/5] Checking all 45 levels...
go test ./...
if errorlevel 1 goto error

echo.
echo [4/5] Preparing output directory...
if not exist "packaging\windows" mkdir "packaging\windows"
if errorlevel 1 goto error

echo.
echo [5/5] Building Windows executable...

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

go build -trimpath -o "packaging\windows\HelwanPenguin.exe" ./cmd/helwan-penguin
if errorlevel 1 goto error

echo.
echo ========================================
echo BUILD SUCCESSFUL
echo Worlds: 9
echo Levels: 45
echo Executable:
echo %CD%\packaging\windows\HelwanPenguin.exe
echo ========================================
echo.
pause
exit /b 0

:error
echo.
echo BUILD FAILED. Check the error messages above.
pause
exit /b 1
