@echo off
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0cursor-host.ps1"
exit /b %ERRORLEVEL%
