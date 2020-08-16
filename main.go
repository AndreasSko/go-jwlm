package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/AndreasSko/go-jwlm/model"
)

func main() {
	var err error

	f, err := os.Create("profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	fmt.Println("Hello world!")
	start := time.Now()
	db := new(model.Database)
	err = db.ImportJWLBackup("UserDataBackup_2020-08-06_Andreas-iPad-2.jwlibrary")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(db.BlockRange))
	fmt.Println(len(db.Bookmark))
	fmt.Println(len(db.Location))
	fmt.Println(len(db.Note))
	fmt.Println(len(db.Tag))
	fmt.Println(len(db.TagMap))
	fmt.Println(len(db.UserMark))

	fmt.Println(db.ExportJWLBackup("test.jwlibrary"))

	duration := time.Since(start)
	fmt.Printf("Ran in %s", duration)
}
