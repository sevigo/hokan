@ECHO off    
if /I %1 == default goto :default
if /I %1 == build goto :build
if /I %1 == run goto :run
if /I %1 == clean goto :clean
if /I %1 == test goto :test

goto :eof ::can be ommited to run the `default` function similarly to makefiles

:default
echo DEFAULT ...
goto :eof

:build
echo BUILD ...
cd cmd\hokan
go build -o ..\..\bin\hokan.exe
cd ..\..
.\bin\hokan.exe
goto :eof

:run
echo RUN ...
cd cmd\hokan
go run main.go inject_server.go inject_store.go inject_watcher.go inject_target.go inject_gui.go wire_gen.go
goto :eof

:test
echo TEST ...
go test -timeout 10s -cover ./...
