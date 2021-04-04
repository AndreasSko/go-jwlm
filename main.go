package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/AndreasSko/go-jwlm/cmd"
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	cmd.Execute()
}
