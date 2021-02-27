package main

import (
	"fmt"
	"syscall/js"

	"github.com/AndreasSko/go-jwlm/wasm"
)

func mergeJs(this js.Value, inputs []js.Value) interface{} {
	leftDbArr := inputs[0]
	rightDbArr := inputs[1]
	mergedDbName := inputs[2].String()
	//mergedDbArr := make([]uint8, leftDbArr.Get("byteLength").Int())

	leftBuf := make([]uint8, leftDbArr.Get("byteLength").Int())
	rightBuf := make([]uint8, rightDbArr.Get("byteLength").Int())

	js.CopyBytesToGo(leftBuf, leftDbArr)
	js.CopyBytesToGo(rightBuf, rightDbArr)

	mergedDb := wasm.Merge(leftBuf, rightBuf, mergedDbName)

	//js.CopyBytesToJS(mergedDbArr, mergedDb)
	fmt.Printf("Merged. Returning %d bytes\n", len(mergedDb))
	mergedJsData := js.Global().Get("Uint8Array").New(len(mergedDb))
	js.CopyBytesToJS(mergedJsData, mergedDb)
	return mergedJsData

}

func registerCallbacks() {
	js.Global().Set("mergeJs", js.FuncOf(mergeJs))
}

func main() {
	//https://blog.twitch.tv/de-de/2019/04/10/go-memory-ballast-how-i-learnt-to-stop-worrying-and-love-the-heap-26c2462549a2/
	ballast := make([]byte, 100<<20) //100MiB
	ballast[0] = 1
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
