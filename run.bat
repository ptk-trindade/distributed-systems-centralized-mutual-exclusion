@echo off
REM r: request mutex r times
REM k: wait k seconds
REM n: run n clientes

set r=3
set k=0
set n=128

REM Start the server

start "" "C:\Users\patrick.c.trindade\Documents\GitHub\distributed-systems-centralized-mutual-exclusion\server\run_server.exe"


REM Start clients
for /l %%i in (1,1,%n%) do (
    start "" "C:\Users\patrick.c.trindade\Documents\GitHub\distributed-systems-centralized-mutual-exclusion\client\run_client.exe" %r% %k%
)
