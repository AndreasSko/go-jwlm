$Env:GOOS="js"; $Env:GOARCH="wasm"; $Env:CGO_ENABLED=0; go build
rm go-jwlm.wasm
mv go-jwlm go-jwlm.wasm -ErrorAction SilentlyContinue