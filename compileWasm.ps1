$Env:GOOS="js"; $Env:GOARCH="wasm"; $Env:CGO_ENABLED=0; go build .\main_libraryhelper.go
rm libhelper.wasm
mv main_libraryhelper libhelper.wasm -ErrorAction SilentlyContinue
cp libhelper.wasm I:\src\libraryhelper\src\assets\