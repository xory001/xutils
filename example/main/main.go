package main

import (
	"github.com/xory001/xutils"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.LUTC)
}

func main() {
	log.Println("xutils example start")
	xutils.SetDebug(true)
	log.Println("xutils example end")
}
