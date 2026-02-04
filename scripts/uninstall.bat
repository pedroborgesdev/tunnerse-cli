@echo off
echo Uninstalling tunnerse CLI and Server for Windows...

set INSTALL_DIR=C:\Program Files\Tunnerse

echo Stopping service if running...
sc query TunnerseServer >nul 2>&1
if %errorlevel% equ 0 (
    echo Stopping TunnerseServer...
    sc stop TunnerseServer
    timeout /t 2 /nobreak >nul
    
    echo Deleting service...
    sc delete TunnerseServer
)

echo Removing binaries...
if exist "%INSTALL_DIR%" (
    rmdir /s /q "%INSTALL_DIR%"
)

echo.
echo Uninstall complete.
echo.
echo Note: PATH variable was not modified automatically.
echo Please remove "%INSTALL_DIR%" from your PATH manually if needed.
echo.
echo Tunnel logs and database were NOT removed.
echo Location: %INSTALL_DIR%\tunnels\
echo.
pause
