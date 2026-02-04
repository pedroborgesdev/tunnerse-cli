@echo off
REM Tunnerse CLI and Server Windows Installer

REM Check for admin rights
openfiles >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Please run this script as Administrator.
    echo Usage: Right-click and 'Run as administrator'
    exit /b 1
)

setlocal
set SCRIPT_DIR=%~dp0
cd /d "%SCRIPT_DIR%"

set BIN_CLI=tunnerse.exe
set BIN_SERVER=tunnerse-server.exe

REM Check if binaries exist
if not exist "%BIN_CLI%" (
    echo Error: %BIN_CLI% not found. Please compile first with: go build -o %BIN_CLI% ..\cmd\cli
    exit /b 1
)
if not exist "%BIN_SERVER%" (
    echo Error: %BIN_SERVER% not found. Please compile first with: go build -o %BIN_SERVER% ..\cmd\server
    exit /b 1
)

REM Install binaries to C:\Program Files\Tunnerse
set INSTALL_DIR=%ProgramFiles%\Tunnerse
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /Y "%BIN_CLI%" "%INSTALL_DIR%\"
copy /Y "%BIN_SERVER%" "%INSTALL_DIR%\"

REM Add to PATH (current session)
setx PATH "%INSTALL_DIR%;%PATH%"

REM Create a basic service using SC
sc create TunnerseServer binPath= '"%INSTALL_DIR%\tunnerse-server.exe"' DisplayName= "Tunnerse Server" start= auto obj= LocalSystem
if %errorlevel% neq 0 (
    echo WARNING: Failed to create service. Run as Administrator!
    echo You can still run tunnerse-server.exe manually.
    goto :skip_service
)

sc description TunnerseServer "Local tunnel management daemon for Tunnerse CLI"
sc failure TunnerseServer reset= 86400 actions= restart/5000/restart/5000/restart/5000
sc start TunnerseServer

:skip_service
echo.
echo Tunnerse CLI and Server installed successfully!
echo Binaries installed to: %INSTALL_DIR%
echo Service installed as: TunnerseServer

echo To manage the service:
echo   sc start TunnerseServer
echo   sc stop TunnerseServer
echo   sc query TunnerseServer
echo   sc delete TunnerseServer

echo Or use Services.msc GUI (services management console)
echo.
echo Use 'tunnerse.exe help' for CLI details.
echo Use 'tunnerse-server.exe' to start the server manually.
echo Or use the Windows service (see above).
endlocal
pause
