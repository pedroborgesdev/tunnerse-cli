@echo off
echo Installing tunnerse CLI and Server for Windows...

cd /d "%~dp0\.."

set BIN_CLI=tunnerse.exe
set BIN_SERVER=tunnerse-server.exe
set BIN_DIR=bin
set INSTALL_DIR=C:\Program Files\Tunnerse

echo Checking for compiled binaries...
if not exist "%BIN_DIR%\%BIN_CLI%" (
    echo ERROR: %BIN_CLI% not found in %BIN_DIR%
    echo Please run build.sh first to compile the binaries.
    exit /b 1
)

if not exist "%BIN_DIR%\%BIN_SERVER%" (
    echo ERROR: %BIN_SERVER% not found in %BIN_DIR%
    echo Please run build.sh first to compile the binaries.
    exit /b 1
)

echo Creating installation directory...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo Installing binaries...
copy /y "%BIN_DIR%\%BIN_CLI%" "%INSTALL_DIR%\"
copy /y "%BIN_DIR%\%BIN_SERVER%" "%INSTALL_DIR%\"

echo Adding to PATH...
setx PATH "%PATH%;%INSTALL_DIR%" /M

echo Creating Windows Service...
sc create TunnerseServer ^
    binPath= "\"%INSTALL_DIR%\%BIN_SERVER%\"" ^
    DisplayName= "Tunnerse Server" ^
    start= auto ^
    obj= LocalSystem

if %errorlevel% neq 0 (
    echo WARNING: Failed to create service. Run as Administrator!
    echo You can still run tunnerse-server.exe manually.
    goto :skip_service
)

echo Configuring service...
sc description TunnerseServer "Local tunnel management daemon for Tunnerse CLI"
sc failure TunnerseServer reset= 86400 actions= restart/5000/restart/5000/restart/5000

echo Starting service...
sc start TunnerseServer

:skip_service
echo.
echo Successfully installed tunnerse CLI and Server.
echo.
echo To manage the service:
echo   Start:   sc start TunnerseServer
echo   Stop:    sc stop TunnerseServer
echo   Status:  sc query TunnerseServer
echo   Delete:  sc delete TunnerseServer
echo.
echo Or use Services.msc GUI (services management console)
echo.
echo Use 'tunnerse help' for CLI details.
echo.
pause
