@echo off 
set CURRENT=%cd%
set LIB_PATH=%CURRENT%\czero\lib
set path=%LIB_PATH%
start /b bin\gaxis.exe attach \\.\pipe\gaxis.ipc

